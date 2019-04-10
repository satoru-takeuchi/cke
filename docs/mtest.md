Multi-host Test (mtest)
=======================

[mtest](../mtest/) directory contains test suites to run integration tests.

Type of Test Suites
-------------------

There are three types of test suites.

1. functions

    This suite tests `ckecli` command and kubernetes workloads deployments which
    related cluster configuration such as etcdbackup.

2. operators

    This suite tests CKE operators user defined resources, labels and taints management.
    If you implement a new operator(Op) in CKE, you have to write test cases in this suite.

Each test suite has an entry point of test as `<suite>/suite_test.go`.

Synopsis
--------

[`Makefile`](../mtest/Makefile) setup virtual machine environment and runs mtest.

* `make setup`

    Install mtest required components.

* `make clean`

    Delete generated files in `output/` directory.

* `make placemat`

    Run `placemat` in background to start virtual machines. To stop placemat, run `sudo pkill placemat`.

* `make test-light`

    Run mtest on a running `placemat`.  This does not control `placemat` by itself.

* `make test`

    Run `make placemat` then run mtest on a running `placemat`. `placemat` is stopped after finish mtest.

* `make bootstrap`

    Create the kubernetes cluster on a running `placemat` using a part of `functions` suite.

Options
-------

### `SUITE`

You can choose the type of test suite by specifying `SUITE` make variable.
The value can be `functions` (default), or `operators`.

`make test` and `make test-light` accepts this variable.

The value of `SUITE` is interpreted as a Go package name.  You can write
a new test suite and specify its package name by `SUITE`.  As a side note,
the forms of `./functions`, and `./operators` are more proper.