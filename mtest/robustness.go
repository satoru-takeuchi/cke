package mtest

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TestStopCP stops 1 control plane for succeeding tests
func TestStopCP() {
	It("should stop CP", func() {
		// stop CKE temporarily to avoid hang-up in SSH session due to node2 shutdown
		stopCKE()

		execAt(node2, "sudo", "systemd-run", "halt", "-f", "-f")
		Eventually(func() error {
			_, err := execAtLocal("ping", "-c", "1", "-W", "1", node2)
			return err
		}).ShouldNot(Succeed())

		execAt(node3, "sudo", "systemctl", "stop", "sshd.socket")

		runCKE(ckeImageURL)
	})
}
