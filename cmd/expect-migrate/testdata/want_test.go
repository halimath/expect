package memfs

import (
	"io/fs"
	"reflect"
	"testing"
	"time"

	"github.com/halimath/expect"
	. "github.com/halimath/fixture"
	"github.com/halimath/fsx"
	"github.com/halimath/expect/is"
)

type memfsFixture struct {
	fs fsx.LinkFS
}

func (f *memfsFixture) BeforeEach(t *testing.T) error {
	f.fs = New()
	return nil
}

func TestMemfs_Mkdir(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("success", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(f.fs.Mkdir("mkdir", 0777))), expect.FailNow(is.NoError(f.fs.Mkdir("mkdir/child", 0777))))

		}).
		Run("noParent", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.Error(f.fs.Mkdir("mkdir/child", 0777), fs.ErrNotExist)))
		}).
		Run("parentNotADirectory", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "not_a_directory", []byte("hello, world"), 0666))), expect.FailNow(is.Error(f.fs.Mkdir("not_a_directory/child", 0777), fs.ErrInvalid)))

		})
}

func TestMemfs_OpenFile(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("success", func(t *testing.T, f *memfsFixture) {
			file, err := f.fs.OpenFile("open_file", fsx.O_RDWR|fsx.O_CREATE, 0644)
			expect.That(t, expect.FailNow(is.NoError(err)))

			l, err := file.Write([]byte("hello, world"))
			expect.That(t, expect.FailNow(is.NoError(err)), expect.FailNow(is.EqualTo(l, len("hello, world"))), expect.FailNow(is.NoError(file.Close())))

			got, err := fs.ReadFile(f.fs, "open_file")
			expect.That(t, expect.FailNow(is.NoError(err)), is.EqualTo(string(got), "hello, world"))

		}).
		Run("notExist", func(t *testing.T, f *memfsFixture) {
			_, err := f.fs.OpenFile("not_found", fsx.O_RDONLY, 0644)
			expect.That(t, expect.FailNow(is.Error(err, fs.ErrNotExist)))
		}).
		Run("parentNotExist", func(t *testing.T, f *memfsFixture) {
			_, err := f.fs.OpenFile("parent_not_found/not_found", fsx.O_RDONLY, 0644)
			expect.That(t, expect.FailNow(is.Error(err, fs.ErrNotExist)))
		}).
		Run("parentNotADirectory", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "not_a_directory", []byte("hello, world"), 0666))))

			_, err := f.fs.OpenFile("not_a_directory/file", fsx.O_CREATE, 0644)
			expect.That(t, expect.FailNow(is.Error(err, fs.ErrInvalid)))
		}).
		Run("parentNotWritable", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(f.fs.Mkdir("dir", 0400))))

			_, err := f.fs.OpenFile("dir/file", fsx.O_WRONLY|fsx.O_CREATE, 0400)
			expect.That(t, is.Error(err, fs.ErrPermission))
		}).
		Run("fileNotWritable", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "file", []byte("hello, world"), 0600))), expect.FailNow(is.NoError(fsx.Chmod(f.fs, "file", 0400))))

			_, err := f.fs.OpenFile("file", fsx.O_WRONLY, 0400)
			expect.That(t, is.Error(err, fs.ErrPermission))
		})
}

func TestMemfs_Open(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("notExist", func(t *testing.T, f *memfsFixture) {
			_, err := f.fs.Open("not_found")
			expect.That(t, expect.FailNow(is.Error(err, fs.ErrNotExist)))
		}).
		Run("success", func(t *testing.T, f *memfsFixture) {
			file, err := f.fs.OpenFile("open", fsx.O_RDWR|fsx.O_CREATE, 0644)
			expect.That(t, expect.FailNow(is.NoError(err)))

			l, err := file.Write([]byte("hello, world"))
			expect.That(t, expect.FailNow(is.NoError(err)), expect.FailNow(is.EqualTo(l, len("hello, world"))), expect.FailNow(is.NoError(file.Close())))

			rf, err := f.fs.Open("open")
			expect.That(t, expect.FailNow(is.NoError(err)))

			buf := make([]byte, len("hello, world"))
			l, err = rf.Read(buf)
			expect.That(t, expect.FailNow(is.NoError(err)), is.EqualTo(l, len("hello, world")), is.EqualTo(string(buf), "hello, world"), expect.FailNow(is.NoError(rf.Close())))

		}).
		Run("dir", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(f.fs.Mkdir("dir", 0777))), expect.FailNow(is.NoError(f.fs.Mkdir("dir/sub_dir", 0777))), expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "dir/sub_file", []byte("hello, world"), 0666))))

			rd, err := f.fs.Open("dir")
			expect.That(t, expect.FailNow(is.NoError(err)))

			info, err := rd.Stat()
			expect.That(t, expect.FailNow(is.NoError(err)), is.EqualTo(info.IsDir(), true))

			readDirFile, ok := rd.(fs.ReadDirFile)
			expect.That(t, expect.FailNow(is.EqualTo(ok, true)))

			entries, err := readDirFile.ReadDir(-1)
			expect.That(t, expect.FailNow(is.NoError(err)), is.DeepEqualTo(entries, []fs.DirEntry{
				&dirEntry{
					name:	"sub_dir",
					info: &fileInfo{
						path:	"dir/sub_dir",
						size:	0,
						mode:	fs.ModeDir | 0777,
					},
				},
				&dirEntry{
					name:	"sub_file",
					info: &fileInfo{
						path:	"dir/sub_file",
						size:	12,
						mode:	0666,
					},
				},
			},
				ExcludeTypes{reflect.TypeOf(time.Now())},
			), expect.FailNow(is.NoError(rd.Close())))

		})
}

func TestMemfs_Remove(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("emptyDir", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(f.fs.Mkdir("dir", 0777))), expect.FailNow(is.NoError(f.fs.Remove("dir"))))

			_, err := fs.Stat(f.fs, "dir")
			expect.That(t, is.Error(err, fs.ErrNotExist))
		}).
		Run("file", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "file", []byte("hello, world"), 0644))), expect.FailNow(is.NoError(f.fs.Remove("file"))))

			_, err := fs.Stat(f.fs, "file")
			expect.That(t, is.Error(err, fs.ErrNotExist))
		}).
		Run("not_exist", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(f.fs.Remove("not_exist"))))
		}).
		Run("parent_not_exist", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.Error(f.fs.Remove("not_exist/sub"), fs.ErrNotExist)))
		}).
		Run("parent_not_a_directory", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "file", []byte("hello, world"), 0644))), expect.FailNow(is.Error(f.fs.Remove("file/sub"), fs.ErrInvalid)))

		})
}

func TestMemfs_SameFile(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("same", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "f1", []byte("hello, world"), 0644))))

			fi1, err := fs.Stat(f.fs, "f1")
			expect.That(t, expect.FailNow(is.NoError(err)))

			fi2, err := fs.Stat(f.fs, "f1")
			expect.That(t, expect.FailNow(is.NoError(err)), is.EqualTo(f.fs.SameFile(fi1, fi2), true))

		}).
		Run("different", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "f1", []byte("hello, world"), 0644))), expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "f2", []byte("hello, world"), 0644))))

			fi1, err := fs.Stat(f.fs, "f1")
			expect.That(t, expect.FailNow(is.NoError(err)))

			fi2, err := fs.Stat(f.fs, "f2")
			expect.That(t, expect.FailNow(is.NoError(err)), is.EqualTo(f.fs.SameFile(fi1, fi2), false))

		})
}

func TestMemfs_Rename(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("old_parent_not_exist", func(t *testing.T, f *memfsFixture) {
			expect.That(t, is.Error(f.fs.Rename("not_exists/file", "file"), fs.ErrNotExist))
		}).
		Run("new_parent_not_exist", func(t *testing.T, f *memfsFixture) {
			expect.That(t, is.Error(f.fs.Rename("file", "not_exists/file"), fs.ErrNotExist))
		}).
		Run("d´directories", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(f.fs.Mkdir("from", 0777))), expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "from/file", []byte("hello, world"), 0644))), expect.FailNow(is.NoError(f.fs.Rename("from", "to"))))

			_, err := fs.Stat(f.fs, "from/file")
			expect.That(t, is.Error(err, fs.ErrNotExist))

			_, err = fs.Stat(f.fs, "to/file")
			expect.That(t, is.NoError(err))
		}).
		Run("file", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "file", []byte("hello, world"), 0644))), expect.FailNow(is.NoError(f.fs.Rename("file", "to"))))

			_, err := fs.Stat(f.fs, "file")
			expect.That(t, is.Error(err, fs.ErrNotExist))

			_, err = fs.Stat(f.fs, "to")
			expect.That(t, is.NoError(err))
		}).
		Run("file_inside_same_dir", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(f.fs.Mkdir("dir", 0777))), expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "dir/from", []byte("hello, world"), 0644))), expect.FailNow(is.NoError(f.fs.Rename("dir/from", "dir/to"))))

			_, err := fs.Stat(f.fs, "dir/from")
			expect.That(t, is.Error(err, fs.ErrNotExist))

			_, err = fs.Stat(f.fs, "dir/to")
			expect.That(t, is.NoError(err))
		}).
		Run("file_between_different_dirs", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(f.fs.Mkdir("from", 0777))), expect.FailNow(is.NoError(f.fs.Mkdir("to", 0777))), expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "from/file", []byte("hello, world"), 0644))), expect.FailNow(is.NoError(f.fs.Rename("from/file", "to/file"))))

			_, err := fs.Stat(f.fs, "from/file")
			expect.That(t, is.Error(err, fs.ErrNotExist))

			_, err = fs.Stat(f.fs, "to/file")
			expect.That(t, is.NoError(err))
		}).
		Run("invalid_source", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "from", []byte("hello, world"), 0644))), is.Error(f.fs.Rename("from/file", "file"), fs.ErrInvalid))

		}).
		Run("invalid_target", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "from", []byte("hello, world"), 0644))), expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "to", []byte("hello, world"), 0644))), is.Error(f.fs.Rename("from", "to/file"), fs.ErrInvalid))

		}).
		Run("source_not_found", func(t *testing.T, f *memfsFixture) {
			expect.That(t, is.Error(f.fs.Rename("from", "to"), fs.ErrNotExist))
		})
}

func TestMemfs_Chmod(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("dir", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(f.fs.Mkdir("dir", 0777))))

			file, err := f.fs.OpenFile("dir", fsx.O_WRONLY, 0)
			expect.That(t, expect.FailNow(is.NoError(err)), expect.FailNow(is.NoError(file.Chmod(0700))), expect.FailNow(is.NoError(file.Close())))

			info, err := fs.Stat(f.fs, "dir")
			expect.That(t, expect.FailNow(is.NoError(err)), is.EqualTo(info.Mode(), fs.ModeDir|0700))

		}).
		Run("file", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "file", []byte("hello, world"), 0666))))

			file, err := f.fs.OpenFile("file", fsx.O_WRONLY, 0)
			expect.That(t, expect.FailNow(is.NoError(err)), expect.FailNow(is.NoError(file.Chmod(0600))), expect.FailNow(is.NoError(file.Close())))

			info, err := fs.Stat(f.fs, "file")
			expect.That(t, expect.FailNow(is.NoError(err)), is.EqualTo(info.Mode(), fs.FileMode(0600)))

		})
}

func TestMemfs_RemoveAll(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("success", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(f.fs.Mkdir("remove_all", 0777))), expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "remove_all/file", []byte("hello, world"), 0644))), expect.FailNow(is.NoError(fsx.RemoveAll(f.fs, "remove_all"))))

			_, err := fs.Stat(f.fs, "remove_all/file")
			expect.That(t, is.Error(err, fs.ErrNotExist))

			_, err = fs.Stat(f.fs, "remove_all")
			expect.That(t, is.Error(err, fs.ErrNotExist))
		})
}

func TestMemfs_ReadFile(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("not_exist", func(t *testing.T, f *memfsFixture) {
			_, err := fs.ReadFile(f.fs, "not_exist")
			expect.That(t, is.Error(err, fs.ErrNotExist))
		}).
		Run("no_file", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(f.fs.Mkdir("dir", 0777))))

			_, err := fs.ReadFile(f.fs, "dir")
			expect.That(t, is.Error(err, ErrIsDirectory))
		}).
		Run("success", func(t *testing.T, f *memfsFixture) {
			expect.That(t, expect.FailNow(is.NoError(fsx.WriteFile(f.fs, "f", []byte("test"), 0644))))

			data, err := fs.ReadFile(f.fs, "f")
			expect.That(t, is.NoError(err), is.DeepEqualTo(data, []byte("test")))

		})
}

func TestMemfs_Symlink(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("success", func(t *testing.T, f *memfsFixture) {
			err := fsx.WriteFile(f.fs, "f", []byte("hello world"), 0666)
			expect.That(t, expect.FailNow(is.NoError(err)), expect.FailNow(is.NoError(f.fs.Symlink("f", "l"))))

			got, err := fs.ReadFile(f.fs, "l")
			expect.That(t, is.NoError(err), is.EqualTo(string(got), "hello world"))

		})
}

func TestMemfs_Link(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("file", func(t *testing.T, f *memfsFixture) {
			err := fsx.WriteFile(f.fs, "f", []byte("hello world"), 0666)
			expect.That(t, expect.FailNow(is.NoError(err)), expect.FailNow(is.NoError(f.fs.Link("f", "l"))))

			got, err := fs.ReadFile(f.fs, "l")
			expect.That(t, is.NoError(err), is.EqualTo(string(got), "hello world"))

		}).
		Run("dir", func(t *testing.T, f *memfsFixture) {
			err := fsx.MkdirAll(f.fs, "dir/child", 0777)
			expect.That(t, expect.FailNow(is.NoError(err)), expect.FailNow(is.NoError(f.fs.Link("dir", "l"))))

			got, err := fs.ReadDir(f.fs, "l")
			expect.That(t, is.NoError(err), is.Len(got, 1), is.EqualTo(got[0].Name(), "child"))

		})
}

func TestMemfs_Readlink(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("symlink", func(t *testing.T, f *memfsFixture) {
			err := fsx.WriteFile(f.fs, "f", []byte("hello world"), 0666)
			expect.That(t, expect.FailNow(is.NoError(err)), expect.FailNow(is.NoError(f.fs.Symlink("f", "l"))))

			got, err := f.fs.Readlink("l")
			expect.That(t, is.NoError(err), is.EqualTo(got, "f"))

		})
}
