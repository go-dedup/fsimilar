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

	// == Test Basic Simhash Function tests
	t.Logf("\n\n== Testing Simhash Basic Functions\n\n")
	// -- sim
	testIt(t, "sim.lstA.sim", "sim -i sim.lstA -d 12 -vv")
	testIt(t, "sim.lstB.sim", "sim -i sim.lstB -d 12 -vv")
	testIt(t, "sim.lstS.sim", "sim -i sim.lstS -d 12 -vv")
	// -- test
	testIt(t, "test1.sim", "sim -i test1.lst -S -d 6 -vv")
	testIt(t, "test2.sim", "sim -i test2.lst -S -d 6 -vv")

	// == Test Simhash Phonetic Function tests
	t.Logf("\n\n== Testing Simhash Phonetic Functions\n\n")
	// -- sim
	testIt(t, "sim.lstA.sim.phonetic", "sim -i sim.lstA -d 12 -p -vv")
	testIt(t, "sim.lstB.sim.phonetic", "sim -i sim.lstB -d 12 -p -vv")
	testIt(t, "sim.lstS.sim.phonetic", "sim -i sim.lstS -d 12 -p -vv")

	// == Test Vector Space Basic Function tests
	t.Logf("\n\n== Testing Vector Space Basic Functions\n\n")
	// -- sim
	testIt(t, "sim.lstA.vec", "vec -i sim.lstA -t 0.7 -vv")
	testIt(t, "sim.lstB.vec", "vec -i sim.lstB -t 0.7 -vv")
	testIt(t, "sim.lstS.vec", "vec -i sim.lstS -t 0.7 -vv")
	// -- test
	testIt(t, "test1.vec", "vec -i test1.lst -S -v")
	testIt(t, "test2.vec", "vec -i test2.lst -S -v")

	// == Test Vector Space Phonetic Function tests
	t.Logf("\n\n== Testing Vector Space Phonetic Functions\n\n")
	// -- sim
	testIt(t, "sim.lstA.vec.phonetic", "vec -i sim.lstA -p -t 0.7 -vv")
	testIt(t, "sim.lstB.vec.phonetic", "vec -i sim.lstB -p -t 0.7 -vv")
	testIt(t, "sim.lstS.vec.phonetic", "vec -i sim.lstS -p -t 0.7 -vv")
	// -- test
	testIt(t, "test1.vec.phonetic", "vec -i test1.lst -S -p -v")
	testIt(t, "test2.vec.phonetic", "vec -i test2.lst -S -p -v")

	// == Test Basic Vector Space Finish Function tests
	t.Logf("\n\n== Testing Vector Space Finish Functions\n\n")
	// -- sim
	testIt(t, "sim.lstA.vec.Finish", "vec -i sim.lstA -t 0.7 -p -F -v")
	testIt(t, "sim.lstB.vec.Finish", "vec -i sim.lstB -t 0.7 -p -F -v")
	testIt(t, "sim.lstS.vec.Finish", "vec -i sim.lstS -t 0.7 -p -F -v")
}
