package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

const (
	cmdTest = "../fsimilar "
	dirTest = "test/"
	extRef  = ".ref" // extension for reference file
	extGot  = ".got" // extension for generated file
)

// testIt runs @cmdEasyGen with @argv and compares the generated
// output for @name with the corresponding @extRef
func testIt(t *testing.T, name string, argv string) {
	var (
		diffOut         bytes.Buffer
		generatedOutput = name + extGot
		cmd             = exec.Command("bash", "-c", cmdTest+argv+" 2>&1")
	)

	t.Logf("Testing %s:\n\t%s%s\n\n", name, cmdTest, argv)

	// open the out file for writing
	outfile, err := os.Create(generatedOutput)
	if err != nil {
		t.Errorf("write error [%s: %s] %s.", name, argv, err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		t.Errorf("start error [%s: %s] %s.", name, argv, err)
	}
	err = cmd.Wait()
	if err != nil {
		t.Errorf("exit error [%s: %s] %s.", name, argv, err)
	}

	cmd = exec.Command("diff", "-U1", name+extRef, generatedOutput)
	cmd.Stdout = &diffOut

	err = cmd.Start()
	if err != nil {
		t.Errorf("start error %s [%s: %s]", err, name, argv)
	}
	err = cmd.Wait()
	if err != nil {
		t.Errorf("cmp error %s [%s: %s]\n%s", err, name, argv, diffOut.String())
	}
	os.Remove(generatedOutput)
}

func TestExec(t *testing.T) {
	os.Chdir(dirTest)

	// == Test Basic Functions
	// -- sim
	testIt(t, "sim.lstA", "-i sim.lstA -d 12 -vv")
	testIt(t, "sim.lstB", "-i sim.lstB -d 12 -vv")
	testIt(t, "sim.lstS", "-i sim.lstS -d 12 -vv")
	// -- test
	testIt(t, "test1", "-i test1.lst -d 6 -vv")
	testIt(t, "test2", "-i test2.lst -d 6 -vv")
}
