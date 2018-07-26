package mtest

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cybozu-go/cke"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"golang.org/x/crypto/ssh"
	yaml "gopkg.in/yaml.v2"
)

const sshTimeout = 3 * time.Minute

var (
	sshClients = make(map[string]*ssh.Client)
)

func sshTo(address string, sshKey ssh.Signer) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: "cybozu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(sshKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	return ssh.Dial("tcp", address+":22", config)
}

func parsePrivateKey() (ssh.Signer, error) {
	f, err := os.Open(os.Getenv("SSH_PRIVKEY"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(data)
}

func prepareSSHClients(addresses ...string) error {
	sshKey, err := parsePrivateKey()
	if err != nil {
		return err
	}

	ch := time.After(sshTimeout)
	for _, a := range addresses {
	RETRY:
		select {
		case <-ch:
			return errors.New("timed out")
		default:
		}
		client, err := sshTo(a, sshKey)
		if err != nil {
			time.Sleep(time.Second)
			goto RETRY
		}
		sshClients[a] = client
	}

	return nil
}

func stopManagementEtcd(client *ssh.Client) error {
	command := "sudo systemctl stop my-etcd.service; sudo rm -rf /home/cybozu/default.etcd"
	sess, err := client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.Run(command)
	return nil
}

func runManagementEtcd(client *ssh.Client) error {
	command := "sudo systemd-run --unit=my-etcd.service /data/etcd --listen-client-urls=http://0.0.0.0:2379 --advertise-client-urls=http://localhost:2379 --data-dir /home/cybozu/default.etcd"
	sess, err := client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	return sess.Run(command)
}

func stopCke() error {
	for _, host := range []string{host1, host2} {
		c := sshClients[host]
		sess, err := c.NewSession()
		if err != nil {
			return err
		}

		sess.Run("sudo systemctl reset-failed cke.service; sudo systemctl stop cke.service")
		sess.Close()
	}
	return nil
}

func runCke() error {
	for _, host := range []string{host1, host2} {
		c := sshClients[host]
		sess, err := c.NewSession()
		if err != nil {
			return err
		}

		err = sess.Run("sudo systemd-run --unit=cke.service --setenv=GOFAIL_HTTP=0.0.0.0:1234 /data/cke -config /etc/cke.yml -interval 10s")
		sess.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func execAt(host string, args ...string) (stdout, stderr []byte, e error) {
	client := sshClients[host]
	sess, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}
	defer sess.Close()

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	sess.Stdout = outBuf
	sess.Stderr = errBuf
	err = sess.Run(strings.Join(args, " "))
	return outBuf.Bytes(), errBuf.Bytes(), err
}

func execSafeAt(host string, args ...string) string {
	stdout, _, err := execAt(host, args...)
	ExpectWithOffset(1, err).To(Succeed())
	return string(stdout)
}

func localTempFile(body string) *os.File {
	f, err := ioutil.TempFile("", "cke-mtest")
	Expect(err).NotTo(HaveOccurred())
	f.WriteString(body)
	f.Close()
	return f
}

func ckecli(args ...string) []byte {
	args = append([]string{"-config", ckeConfigPath}, args...)
	command := exec.Command(ckecliPath, args...)
	stdout := new(bytes.Buffer)
	session, err := gexec.Start(command, stdout, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gexec.Exit(0))
	return stdout.Bytes()
}

func getClusterStatus() (*cke.ClusterStatus, error) {
	controller := cke.NewController(nil, 0)

	f, err := os.Open(ckeClusterPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cluster cke.Cluster
	err = yaml.NewDecoder(f).Decode(&cluster)
	if err != nil {
		return nil, err
	}

	return controller.GetClusterStatus(context.Background(), &cluster)
}
