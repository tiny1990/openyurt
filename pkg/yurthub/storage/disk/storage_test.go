package disk

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	tempDir = "kubelet/default/pods"
	tempKey = "kubelet/default/pods/test-pod"
)

func TestCreate(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	err = s.Create(tempKey, []byte("test-pod"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s", err, tempKey)
	}

	createdFile := filepath.Join(cacheBaseDir, tempKey)
	if fi, err := os.Stat(createdFile); err != nil {
		t.Errorf("Got error %v, wanted file %q to be there", err, createdFile)
	} else if !fi.Mode().IsRegular() {
		t.Errorf("Got %q not a regular file", createdFile)
	}

	b, err := ioutil.ReadFile(createdFile)
	if err != nil {
		t.Errorf("Got error %v, unable read regular file %q", err, createdFile)
	} else if !bytes.Equal(b, []byte("test-pod")) {
		t.Errorf("Wanted string: test-pod but got %s", string(b))
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestCreateFileExist(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	err = s.Create(tempKey, []byte("test-pod1"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s", err, tempKey)
	}

	err = s.Create(tempKey, []byte("test-pod2"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s witch contents test-pod2", err, tempKey)
	}

	createdFile := filepath.Join(cacheBaseDir, tempKey)
	if fi, err := os.Stat(createdFile); err != nil {
		t.Errorf("Got error %v, wanted file %q to be there", err, createdFile)
	} else if !fi.Mode().IsRegular() {
		t.Errorf("Got %q not a regular file", createdFile)
	}

	b, err := ioutil.ReadFile(createdFile)
	if err != nil {
		t.Errorf("Got error %v, unable read regular file %q", err, createdFile)
	} else if !bytes.Equal(b, []byte("test-pod2")) {
		t.Errorf("Wanted string: test-pod2 but got %s", string(b))
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestCreateDirExist(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	createdFile := filepath.Join(cacheBaseDir, tempKey)
	dir, _ := filepath.Split(createdFile)
	if err = os.MkdirAll(dir, 0755); err != nil {
		t.Errorf("Got error %v, unable make dir %s", err, dir)
	}

	err = s.Create(tempKey, []byte("test-pod"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s", err, tempKey)
	}

	if fi, err := os.Stat(createdFile); err != nil {
		t.Errorf("Got error %v, wanted file %q to be there", err, createdFile)
	} else if !fi.Mode().IsRegular() {
		t.Errorf("Got %q not a regular file", createdFile)
	}

	b, err := ioutil.ReadFile(createdFile)
	if err != nil {
		t.Errorf("Got error %v, unable read regular file %q", err, createdFile)
	} else if !bytes.Equal(b, []byte("test-pod")) {
		t.Errorf("Wanted string: test-pod but got %s", string(b))
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestDelete(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	err = s.Create(tempKey, []byte("test-pod"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s", err, tempKey)
	}

	createdFile := filepath.Join(cacheBaseDir, tempKey)
	if fi, err := os.Stat(createdFile); err != nil {
		t.Errorf("Got error %v, wanted file %q to be there", err, createdFile)
	} else if !fi.Mode().IsRegular() {
		t.Errorf("Got %q not a regular file", createdFile)
	}

	err = s.Delete(tempKey)
	if err != nil {
		t.Errorf("Got error %v, unable delete key %q", err, tempKey)
	}

	if _, err := os.Stat(createdFile); err == nil || !os.IsNotExist(err) {
		t.Errorf("want %q is deleted, but it still exist", createdFile)
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestDeleteFileNotExist(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	createdFile := filepath.Join(cacheBaseDir, tempKey)
	err = s.Delete(tempKey)
	if err != nil {
		t.Errorf("Got error %v, delete not exist file(%q) returned error", err, createdFile)
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestDeleteDir(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	err = s.Create(tempKey, []byte("test-pod"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s", err, tempKey)
	}

	err = s.Delete(tempDir)
	if err != nil {
		t.Errorf("Got error %v, unable delete dir key %q", err, tempDir)
	}

	createdFile := filepath.Join(cacheBaseDir, tempKey)
	if fi, err := os.Stat(createdFile); err != nil {
		t.Errorf("Got error %v, wanted file %q to be there", err, createdFile)
	} else if !fi.Mode().IsRegular() {
		t.Errorf("Got %q not a regular file", createdFile)
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestGet(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	err = s.Create(tempKey, []byte("test-pod"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s", err, tempKey)
	}

	b, err := s.Get(tempKey)
	if err != nil {
		t.Errorf("Got error %v, get key %q", err, tempKey)
	} else if !bytes.Equal(b, []byte("test-pod")) {
		t.Errorf("Wanted string: test-pod but got %s", string(b))
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestGetFileNotExist(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	b, err := s.Get(tempKey)
	if err != nil {
		t.Errorf("Got error %v, get key %q", err, tempKey)
	} else if len(b) != 0 {
		t.Errorf("Wanted empty string got %s", string(b))
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestGetNotRegularFile(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	err = s.Create(tempKey, []byte("test-pod"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s", err, tempKey)
	}

	_, err = s.Get(tempDir)
	if err == nil {
		t.Errorf("Got not error for dir key %q", tempDir)
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestListKeys(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	tempKeys := make([]string, 5)
	for i := 0; i < 5; i++ {
		tempKeys[i] = fmt.Sprintf("%s-%d", tempKey, i)
		err = s.Create(tempKeys[i], []byte("test-pod"))
		if err != nil {
			t.Errorf("Got error %v, wanted successful create %s", err, tempKeys[i])
		}
	}

	keys, err := s.ListKeys(tempDir)
	if err != nil {
		t.Errorf("Got error %v, unable list keys for %s", err, tempDir)
	}

	if len(tempKeys) != len(keys) {
		t.Errorf("expect %d keys, but got %d keys", len(tempKeys), len(keys))
	}

	for _, key := range tempKeys {
		found := false
		for _, cachedKey := range keys {
			if key == cachedKey {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("key %s is not found by list keys", key)
		}
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestListKeysForEmptyDir(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	keys, err := s.ListKeys(tempDir)
	if err != nil {
		t.Errorf("Got error %v, unable list keys for empty dir %s", err, tempDir)
	}

	if len(keys) != 0 {
		t.Errorf("expect 0 key, but got %d keys", len(keys))
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestListKeysForRegularFile(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	err = s.Create(tempKey, []byte("test-pod"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s", err, tempKey)
	}

	keys, err := s.ListKeys(tempKey)
	if err != nil {
		t.Errorf("Got error %v, unable list keys for empty dir %s", err, tempDir)
	}

	if len(keys) != 1 {
		t.Errorf("listKeys: expect 1 key, but got %d keys", len(keys))
	}

	if keys[0] != tempKey {
		t.Errorf("listKeys: expect %s key, but got %s key", tempKey, keys[0])
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestList(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	tempContents := make([]string, 5)
	for i := 0; i < 5; i++ {
		tempContents[i] = fmt.Sprintf("test-pod-%d", i)
		err = s.Create(fmt.Sprintf("%s-%d", tempKey, i), []byte(tempContents[i]))
		if err != nil {
			t.Errorf("Got error %v, wanted successful create %s", err, fmt.Sprintf("%s-%d", tempKey, i))
		}
	}

	contents, err := s.List(tempDir)
	if err != nil {
		t.Errorf("Got error %v, unable list for %s", err, tempDir)
	}

	if len(tempContents) != len(contents) {
		t.Errorf("expect %d number of contents, but got %d number of contents", len(tempContents), len(contents))
	}

	for _, content := range tempContents {
		found := false
		for _, cachedContent := range contents {
			if content == string(cachedContent) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("content %s is not found by list", content)
		}
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestListEmptyDir(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	contents, err := s.List(tempDir)
	if err != nil {
		t.Errorf("Got error %v, unable list for %s", err, tempDir)
	}

	if len(contents) != 0 {
		t.Errorf("expect no contents, but got %d number of contents", len(contents))
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestListSpecifiedFile(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	err = s.Create(tempKey, []byte("test-pod"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s", err, tempKey)
	}

	contents, err := s.List(tempKey)
	if err != nil {
		t.Errorf("Got error %v, unable list for %s", err, tempKey)
	}

	if len(contents) != 1 {
		t.Errorf("expect 1 contents, but got %d number of contents", len(contents))
	}

	if string(contents[0]) != "test-pod" {
		t.Errorf("expect content: test-pod, but got content: %s", contents[0])
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestUpdate(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	err = s.Create(tempKey, []byte("test-pod"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s", err, tempKey)
	}

	err = s.Update(tempKey, []byte("test-pod1"))
	if err != nil {
		t.Errorf("Got error %v, unable update key %s", err, tempKey)
	}

	createdFile := filepath.Join(cacheBaseDir, tempKey)
	if fi, err := os.Stat(createdFile); err != nil {
		t.Errorf("Got error %v, wanted file %q to be there", err, createdFile)
	} else if !fi.Mode().IsRegular() {
		t.Errorf("Got %q not a regular file", createdFile)
	}

	b, err := ioutil.ReadFile(createdFile)
	if err != nil {
		t.Errorf("Got error %v, unable read regular file %q", err, createdFile)
	} else if !bytes.Equal(b, []byte("test-pod1")) {
		t.Errorf("Wanted string: test-pod1 but got %s", string(b))
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}

func TestUpdateEmptyString(t *testing.T) {
	s, err := NewDiskStorage()
	if err != nil {
		t.Fatalf("unable to new disk storage, %v", err)
	}

	err = s.Create(tempKey, []byte("test-pod"))
	if err != nil {
		t.Errorf("Got error %v, wanted successful create %s", err, tempKey)
	}

	err = s.Update(tempKey, []byte(""))
	if err != nil {
		t.Errorf("Got error %v, unable update key %s", err, tempKey)
	}

	createdFile := filepath.Join(cacheBaseDir, tempKey)
	if fi, err := os.Stat(createdFile); err != nil {
		t.Errorf("Got error %v, wanted file %q to be there", err, createdFile)
	} else if !fi.Mode().IsRegular() {
		t.Errorf("Got %q not a regular file", createdFile)
	}

	b, err := ioutil.ReadFile(createdFile)
	if err != nil {
		t.Errorf("Got error %v, unable read regular file %q", err, createdFile)
	} else if len(b) == 0 {
		t.Errorf("Wanted string: empty string but got %s", string(b))
	}

	if err = os.RemoveAll(cacheBaseDir); err != nil {
		t.Errorf("Got error %v, unable remove path %s", err, cacheBaseDir)
	}
}
