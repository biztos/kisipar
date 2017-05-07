// filesystemprovider_test.go -- tests for filesystem provider
// --------------------------

package kisipar_test

import (
	// Standard:
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	// Helpful:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/kisipar"
)

func Test_NewFileSystemProviderConfig_UnknownKey(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{"fremd": true}
	_, err := kisipar.NewFileSystemProviderConfig(d)
	if assert.Error(err) {
		assert.Equal("Unexpected FileSystemProviderConfig key: fremd",
			err.Error(), "error useful")
	}

}

func Test_NewFileSystemProviderConfig_ContentDirNotString(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{"ContentDir": true}
	_, err := kisipar.NewFileSystemProviderConfig(d)
	if assert.Error(err) {
		assert.Equal("ContentDir must be a string, not bool.",
			err.Error(), "error useful")
	}

}

func Test_NewFileSystemProviderConfig_TemplateDirNotString(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{"TemplateDir": true}
	_, err := kisipar.NewFileSystemProviderConfig(d)
	if assert.Error(err) {
		assert.Equal("TemplateDir must be a string, not bool.",
			err.Error(), "error useful")
	}

}

func Test_NewFileSystemProviderConfig_TemplateThemeNotString(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{"TemplateTheme": true}
	_, err := kisipar.NewFileSystemProviderConfig(d)
	if assert.Error(err) {
		assert.Equal("TemplateTheme must be a string, not bool.",
			err.Error(), "error useful")
	}

}

func Test_NewFileSystemProviderConfig_ExcludeWrongType(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{"Exclude": true}
	_, err := kisipar.NewFileSystemProviderConfig(d)
	if assert.Error(err) {
		assert.Equal("Exclude is neither a *regexp.Regexp nor a string, but a bool.",
			err.Error(), "error useful")
	}

}

func Test_NewFileSystemProviderConfig_ExcludeBadRegex(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{"Exclude": "^foo[a"}
	_, err := kisipar.NewFileSystemProviderConfig(d)
	if assert.Error(err) {
		assert.Regexp("^Exclude is not a valid regexp string",
			err.Error(), "error useful")
	}

}

func Test_NewFileSystemProviderConfig_IncludeWrongType(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{"Include": true}
	_, err := kisipar.NewFileSystemProviderConfig(d)
	if assert.Error(err) {
		assert.Equal("Include is neither a *regexp.Regexp nor a string, but a bool.",
			err.Error(), "error useful")
	}

}

func Test_NewFileSystemProviderConfig_IncludeBadRegex(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{"Include": "^foo[a"}
	_, err := kisipar.NewFileSystemProviderConfig(d)
	if assert.Error(err) {
		assert.Regexp("^Include is not a valid regexp string",
			err.Error(), "error useful")
	}

}

func Test_NewFileSystemProviderConfig_AllowMetaErrorsNotBool(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{"AllowMetaErrors": 123}
	_, err := kisipar.NewFileSystemProviderConfig(d)
	if assert.Error(err) {
		assert.Equal("AllowMetaErrors must be a bool, not int.",
			err.Error(), "error useful")
	}

}

func Test_NewFileSystemProviderConfig_AutoRefreshNotBool(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{"AutoRefresh": 123}
	_, err := kisipar.NewFileSystemProviderConfig(d)
	if assert.Error(err) {
		assert.Equal("AutoRefresh must be a bool, not int.",
			err.Error(), "error useful")
	}

}

func Test_NewFileSystemProviderConfig_SuccessWithEmpty(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{}
	cfg, err := kisipar.NewFileSystemProviderConfig(d)
	if !assert.Nil(err) {
		assert.FailNow(err.Error())
	}

	exp := &kisipar.FileSystemProviderConfig{}
	assert.Equal(exp, cfg, "got empty config")

}

func Test_NewFileSystemProviderConfig_SuccessWithZeroValues(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{
		"ContentDir":      "",
		"TemplateDir":     "",
		"TemplateTheme":   "",
		"Exclude":         nil,
		"Include":         nil,
		"AllowMetaErrors": false,
		"AutoRefresh":     false,
	}
	cfg, err := kisipar.NewFileSystemProviderConfig(d)
	if !assert.Nil(err) {
		assert.FailNow(err.Error())
	}

	exp := &kisipar.FileSystemProviderConfig{}
	assert.Equal(exp, cfg, "got empty config")

}

func Test_NewFileSystemProviderConfig_SuccessWithNonzeroValues(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{
		"ContentDir":      "cdir/",
		"TemplateDir":     "tdir/",
		"TemplateTheme":   "happy",
		"Exclude":         regexp.MustCompile("ex"),
		"Include":         regexp.MustCompile("in"),
		"AllowMetaErrors": true,
		"AutoRefresh":     true,
	}
	cfg, err := kisipar.NewFileSystemProviderConfig(d)
	if !assert.Nil(err) {
		assert.FailNow(err.Error())
	}

	exp := &kisipar.FileSystemProviderConfig{
		ContentDir:      "cdir/",
		TemplateDir:     "tdir/",
		TemplateTheme:   "happy",
		Exclude:         regexp.MustCompile("ex"),
		Include:         regexp.MustCompile("in"),
		AllowMetaErrors: true,
		AutoRefresh:     true,
	}
	assert.Equal(exp, cfg, "got expected config")

}

func Test_NewFileSystemProviderConfig_SuccessWithStringRegexps(t *testing.T) {

	assert := assert.New(t)

	d := map[string]interface{}{
		"Exclude": "ex",
		"Include": "in",
	}
	cfg, err := kisipar.NewFileSystemProviderConfig(d)
	if !assert.Nil(err) {
		assert.FailNow(err.Error())
	}

	exp := &kisipar.FileSystemProviderConfig{
		Exclude: regexp.MustCompile("ex"),
		Include: regexp.MustCompile("in"),
	}
	assert.Equal(exp, cfg, "got expected config")

}

func Test_NewFileSystemProvider(t *testing.T) {

	assert := assert.New(t)

	fsp := kisipar.NewFileSystemProvider(&kisipar.FileSystemProviderConfig{
		ContentDir: "/any/where",
	})

	// Note that the the update time is initialized.
	assert.Regexp(
		"^<FileSystemProvider with 0 items at /any/where, updated .*>$",
		fsp.String(), "stringifies as expected")

}

func Test_FileSystemProvider_LoadTemplates_NoTemplateDir(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadTemplates()
	assert.Nil(err, "no error")
	assert.Nil(fsp.Template(), "no master template set")

}

func Test_FileSystemProvider_LoadTemplates_NoSuchDir(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{TemplateDir: "nosuchdir"}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadTemplates()
	if assert.Error(err, "got error") {
		assert.Regexp(".*no such file or directory$", err.Error())
	}
}

func Test_FileSystemProvider_LoadTemplates_DirNotDir(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{TemplateDir: "README.md"}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadTemplates()
	if assert.Error(err, "got error") {
		assert.Equal("Not a directory: README.md", err.Error())
	}
}

func Test_FileSystemProvider_LoadTemplates_TemplateError(t *testing.T) {

	assert := assert.New(t)

	dir := filepath.Join("testdata", "fsp-bad-templates")
	config := &kisipar.FileSystemProviderConfig{TemplateDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadTemplates()
	if assert.Error(err, "got error") {
		assert.Regexp("^Error walking .* Template", err.Error())
	}
}

// Sort of a pain in the ass edge case but I hit it for real while debugging
// so (alas) it's worth testing for.
func Test_FileSystemProvider_LoadTemplates_InnerSymlinkError(t *testing.T) {

	assert := assert.New(t)

	// Top dir to hold the goods.
	dir, err := ioutil.TempDir("", "kisipar-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Inner file which is the link target.
	fn := filepath.Join(dir, "target.html")
	if err = ioutil.WriteFile(fn, []byte("hello"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Link to it.
	ln := filepath.Join(dir, "linked.html")
	if err = os.Symlink(fn, ln); err != nil {
		t.Fatal(err)
	}

	// Remove the original but leave the link...
	os.Remove(fn)

	// And we should get an error loading the templates:
	config := &kisipar.FileSystemProviderConfig{TemplateDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err = fsp.LoadTemplates()
	if assert.Error(err, "got error") {
		assert.Regexp("^Error walking .*linked.html", err.Error())
	}

}

func Test_FileSystemProvider_LoadTemplates_InnerSymlinkIsDir(t *testing.T) {

	assert := assert.New(t)

	// Top dir to hold the goods.
	dir, err := ioutil.TempDir("", "kisipar-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Inner *dir* which is the link target.
	tdir, err := ioutil.TempDir(dir, "target.html")
	if err != nil {
		t.Fatal(err)
	}

	// Link to it (the link will be a file).
	ln := filepath.Join(dir, "linked.html")
	if err = os.Symlink(tdir, ln); err != nil {
		t.Fatal(err)
	}

	// And we should get no error, but exercise the symlink dir skip logic:
	config := &kisipar.FileSystemProviderConfig{TemplateDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err = fsp.LoadTemplates()
	if !assert.Nil(err, "no error") {
		t.Log(err)
	}

}

func Test_FileSystemProvider_LoadTemplates_FileReadError(t *testing.T) {

	assert := assert.New(t)

	// Top dir to hold the goods.
	dir, err := ioutil.TempDir("", "kisipar-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Inner file which we will fail to read.
	fn := filepath.Join(dir, "tilos.html")
	if err = ioutil.WriteFile(fn, []byte("hello"), os.FileMode(0000)); err != nil {
		t.Fatal(err)
	}

	// And we should get a file read error walking the dir:
	config := &kisipar.FileSystemProviderConfig{TemplateDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err = fsp.LoadTemplates()
	if assert.Error(err, "got error") {
		assert.Regexp("^Error walking .*tilos.html", err.Error())
	}

}

func Test_FileSystemProvider_LoadTemplates_Success(t *testing.T) {

	assert := assert.New(t)

	dir := filepath.Join("testdata", "fsp-templates")
	config := &kisipar.FileSystemProviderConfig{TemplateDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadTemplates()
	if !assert.Nil(err, "no error") {
		t.Log(err)
	}
}

func Test_FileSystemProvider_LoadInternalTemplates_ThemeError(t *testing.T) {

	assert := assert.New(t)

	// Do all of them.
	config := &kisipar.FileSystemProviderConfig{TemplateTheme: "nonesuch"}
	fsp := kisipar.NewFileSystemProvider(config)
	err := fsp.LoadInternalTemplates()
	if assert.Error(err) {
		assert.Equal("No templates available for theme nonesuch.", err.Error(),
			"error as expected")
	}

}

func Test_FileSystemProvider_LoadInternalTemplates_Success(t *testing.T) {

	assert := assert.New(t)

	// Do all of them.
	themes := kisipar.TemplateThemes()
	if len(themes) < 1 {
		panic("no themes from TemplateThemes")
	}
	for _, theme := range themes {
		config := &kisipar.FileSystemProviderConfig{TemplateTheme: theme}
		fsp := kisipar.NewFileSystemProvider(config)
		err := fsp.LoadInternalTemplates()
		if !assert.Nil(err, "no error") {
			t.Log(err)
		}
	}

}

func Test_FileSystemProvider_LoadInternalTemplates_SuccessDefault(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{}
	fsp := kisipar.NewFileSystemProvider(config)
	err := fsp.LoadInternalTemplates()
	if !assert.Nil(err, "no error") {
		t.Log(err)
	}

}

func Test_FileSystemProvider_LoadContent_NoContentDir(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{ContentDir: ""}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	if !assert.Nil(err, "no error") {
		t.Logf("Error: %s", err.Error())
	}
}

func Test_FileSystemProvider_LoadContent_NoSuchDir(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{ContentDir: "nosuchdir"}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	if assert.Error(err, "got error") {
		assert.Regexp(".*no such file or directory$", err.Error())
	}
}

func Test_FileSystemProvider_LoadContent_DirNotDir(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{ContentDir: "README.md"}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	if assert.Error(err, "got error") {
		assert.Equal("Not a directory: README.md", err.Error())
	}
}

// Same story as for templates:
func Test_FileSystemProvider_LoadContent_InnerSymlinkError(t *testing.T) {

	assert := assert.New(t)

	// Top dir to hold the goods.
	dir, err := ioutil.TempDir("", "kisipar-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Inner file which is the link target.
	fn := filepath.Join(dir, "target.html")
	if err = ioutil.WriteFile(fn, []byte("hello"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Link to it.
	ln := filepath.Join(dir, "linked.html")
	if err = os.Symlink(fn, ln); err != nil {
		t.Fatal(err)
	}

	// Remove the original but leave the link...
	os.Remove(fn)

	// And we should get an error loading the content:
	config := &kisipar.FileSystemProviderConfig{ContentDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err = fsp.LoadContent()
	if assert.Error(err, "got error") {
		assert.Regexp("^Error walking .*linked.html", err.Error())
	}

}

func Test_FileSystemProvider_LoadContent_InnerSymlinkIsDir(t *testing.T) {

	assert := assert.New(t)

	// Top dir to hold the goods.
	dir, err := ioutil.TempDir("", "kisipar-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Inner *dir* which is the link target.
	tdir, err := ioutil.TempDir(dir, "target.html")
	if err != nil {
		t.Fatal(err)
	}

	// Link to it (the link will be a file).
	ln := filepath.Join(dir, "linked.html")
	if err = os.Symlink(tdir, ln); err != nil {
		t.Fatal(err)
	}

	// And we should get no error, but exercise the symlink dir skip logic:
	config := &kisipar.FileSystemProviderConfig{ContentDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err = fsp.LoadContent()
	if !assert.Nil(err, "no error") {
		t.Log(err)
	}

}

func Test_FileSystemProvider_LoadContent_FileReadError(t *testing.T) {

	assert := assert.New(t)

	// Top dir to hold the goods.
	dir, err := ioutil.TempDir("", "kisipar-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Inner file which we will fail to read.
	fn := filepath.Join(dir, "tilos.md")
	if err = ioutil.WriteFile(fn, []byte("hello"), os.FileMode(0000)); err != nil {
		t.Fatal(err)
	}

	// And we should get a file read error walking the dir:
	config := &kisipar.FileSystemProviderConfig{ContentDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err = fsp.LoadContent()
	if assert.Error(err, "got error") {
		assert.Regexp("^Error walking .*tilos.md", err.Error())
	}

}

func Test_FileSystemProvider_LoadContent_PageMetaError(t *testing.T) {

	assert := assert.New(t)

	dir := filepath.Join("testdata", "fsp-bad-content")

	// When Strict is set, we get an error:
	config := &kisipar.FileSystemProviderConfig{ContentDir: dir}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	if assert.Error(err, "got error") {
		assert.Regexp("^Error walking .*bad-yaml.md.*yaml", err.Error())
	}

}

func Test_LoadFileSystemProvider_PageMetaError_AllowMetaErrors(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{
		ContentDir:      filepath.Join("testdata", "fsp-bad-content"),
		AllowMetaErrors: true,
	}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	if !assert.Nil(err, "no error") {
		t.Log(err)
	}

	exp := []string{
		"/bad-yaml",
	}
	assert.Equal(exp, fsp.Paths(), "paths as expected")
}

func Test_FileSystemProvider_LoadContent_Success(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{
		ContentDir: filepath.Join("testdata", "fsp-content"),
	}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	assert.Nil(err, "no error")

	exp := []string{
		"/dupe",
		"/foo",
		"/index",
		"/other",
		"/other.txt",
		"/foo/bar",
		"/foo/s.js",
		"/foo/bar/baz",
		"/foo/bother/data.json",
		"/foo/bother/boo/bam",
	}
	assert.Equal(exp, fsp.Paths(), "paths as expected")

	// uh-oh, we had some trouble loading the right types, so...
	ptypes := map[string]string{
		"/dupe":                 "*kisipar.StandardPage",
		"/foo":                  "*kisipar.StandardPage",
		"/index":                "*kisipar.StandardPage",
		"/other":                "*kisipar.StandardPage",
		"/other.txt":            "*kisipar.StandardFile",
		"/foo/bar":              "*kisipar.StandardPage",
		"/foo/s.js":             "*kisipar.StandardFile",
		"/foo/bar/baz":          "*kisipar.StandardPage",
		"/foo/bother/data.json": "*kisipar.StandardFile",
		"/foo/bother/boo/bam":   "*kisipar.StandardPage",
	}
	for path, pt := range ptypes {
		p, err := fsp.Get(path)
		if assert.Nil(err) {
			assert.Equal(pt, fmt.Sprintf("%T", p), "%s is a %s", path, pt)
		}
	}
}

func Test_FileSystemProvider_LoadContent_Success_Exclude(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{
		ContentDir: filepath.Join("testdata", "fsp-content"),
		Exclude:    regexp.MustCompile("^foo/b|other.md"),
	}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	assert.Nil(err, "no error")

	// Excluded items are commented out:
	exp := []string{
		"/dupe",
		"/foo",
		"/index",
		// "/other",
		"/other.txt",
		// "/foo/bar",
		"/foo/s.js",
		// "/foo/bar/baz",
		// "/foo/bother/data.json",
		// "/foo/bother/boo/bam",
	}
	assert.Equal(exp, fsp.Paths(), "paths as expected")
}

func Test_FileSystemProvider_LoadContent_Success_Include(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{
		ContentDir: filepath.Join("testdata", "fsp-content"),
		Include:    regexp.MustCompile("[.](json|txt)$"),
	}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	assert.Nil(err, "no error")

	// Excluded items are commented out:
	exp := []string{
		// "/dupe",
		// "/foo",
		// "/index",
		// "/other",
		"/other.txt",
		// "/foo/bar",
		// "/foo/s.js",
		// "/foo/bar/baz",
		"/foo/bother/data.json",
		// "/foo/bother/boo/bam",
	}
	assert.Equal(exp, fsp.Paths(), "paths as expected")
}

func Test_FileSystemProvider_LoadContent_Success_ExcludeInclude(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{
		ContentDir: filepath.Join("testdata", "fsp-content"),
		Exclude:    regexp.MustCompile("^other"),
		Include:    regexp.MustCompile("[.](json|txt)$"),
	}

	fsp := kisipar.NewFileSystemProvider(config)

	err := fsp.LoadContent()
	assert.Nil(err, "no error")

	// Excluded items are commented out:
	exp := []string{
		// "/dupe",
		// "/foo",
		// "/index",
		// "/other",
		// "/other.txt",
		// "/foo/bar",
		// "/foo/s.js",
		// "/foo/bar/baz",
		"/foo/bother/data.json",
		// "/foo/bother/boo/bam",
	}
	assert.Equal(exp, fsp.Paths(), "paths as expected")
}

func Test_LoadFileSystemProvider_TemplateError(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{
		ContentDir:  filepath.Join("testdata", "fsp-content"),
		TemplateDir: filepath.Join("testdata", "fsp-bad-templates"),
	}

	_, err := kisipar.LoadFileSystemProvider(config)

	if assert.Error(err) {
		assert.Regexp("^Error walking .*broken.html", err.Error())
	}
}

func Test_LoadFileSystemProvider_ContentError(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{
		ContentDir:  filepath.Join("testdata", "fsp-bad-content"),
		TemplateDir: filepath.Join("testdata", "fsp-templates"),
	}

	_, err := kisipar.LoadFileSystemProvider(config)

	if assert.Error(err) {
		assert.Regexp("^Error walking .*bad-yaml.md.*yaml", err.Error())
	}
}

func Test_LoadFileSystemProvider_TemplateThemeError(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{
		ContentDir:    filepath.Join("testdata", "fsp-bad-content"),
		TemplateTheme: "nonesuch",
	}

	_, err := kisipar.LoadFileSystemProvider(config)

	if assert.Error(err) {
		assert.Equal("No templates available for theme nonesuch.", err.Error(),
			"error as expected")
	}
}

func Test_LoadFileSystemProvider_Success(t *testing.T) {

	assert := assert.New(t)

	config := &kisipar.FileSystemProviderConfig{
		ContentDir:  filepath.Join("testdata", "fsp-content"),
		TemplateDir: filepath.Join("testdata", "fsp-templates"),
	}

	fsp, err := kisipar.LoadFileSystemProvider(config)

	assert.Nil(err, "no error")

	exp := []string{
		"/dupe",
		"/foo",
		"/index",
		"/other",
		"/other.txt",
		"/foo/bar",
		"/foo/s.js",
		"/foo/bar/baz",
		"/foo/bother/data.json",
		"/foo/bother/boo/bam",
	}
	assert.Equal(exp, fsp.Paths(), "paths as expected")
}
