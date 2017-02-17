/*
 * Minio Client (C) 2014, 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"net/http/httptest"

	"github.com/hashicorp/go-version"
	"github.com/minio/cli"
	. "gopkg.in/check.v1"
)

var customConfigDir string

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

var server *httptest.Server
var app *cli.App

func (s *TestSuite) SetUpSuite(c *C) {
}

func (s *TestSuite) TearDownSuite(c *C) {
}

func (s *TestSuite) TestValidPERMS(c *C) {
	perms := accessPerms("none")
	c.Assert(perms.isValidAccessPERM(), Equals, true)
	c.Assert(string(perms), Equals, "none")
	perms = accessPerms("public")
	c.Assert(perms.isValidAccessPERM(), Equals, true)
	c.Assert(string(perms), Equals, "public")
	perms = accessPerms("download")
	c.Assert(perms.isValidAccessPERM(), Equals, true)
	c.Assert(string(perms), Equals, "download")
	perms = accessPerms("upload")
	c.Assert(perms.isValidAccessPERM(), Equals, true)
	c.Assert(string(perms), Equals, "upload")
}

func (s *TestSuite) TestInvalidPERMS(c *C) {
	perms := accessPerms("invalid")
	c.Assert(perms.isValidAccessPERM(), Equals, false)
}

func (s *TestSuite) TestGetMiniocConfigDir(c *C) {
	dir, err := getMiniocConfigDir()
	c.Assert(err, IsNil)
	c.Assert(dir, Not(Equals), "")
	c.Assert(mustGetMiniocConfigDir(), Equals, dir)
}

func (s *TestSuite) TestGetMiniocConfigPath(c *C) {
	dir, err := getMiniocConfigPath()
	c.Assert(err, IsNil)
	switch runtime.GOOS {
	case "linux", "freebsd", "darwin", "solaris":
		c.Assert(dir, Equals, filepath.Join(mustGetMiniocConfigDir(), "config.json"))
	case "windows":
		c.Assert(dir, Equals, filepath.Join(mustGetMiniocConfigDir(), "config.json"))
	default:
		c.Fail()
	}
	c.Assert(mustGetMiniocConfigPath(), Equals, dir)
}

func (s *TestSuite) TestIsvalidAliasName(c *C) {
	c.Check(isValidAlias("helloWorld0"), Equals, true)
	c.Check(isValidAlias("h0SFD2k24Fdsa"), Equals, true)
	c.Check(isValidAlias("fdslka-4"), Equals, true)
	c.Check(isValidAlias("fdslka-"), Equals, true)
	c.Check(isValidAlias("helloWorld$"), Equals, false)
	c.Check(isValidAlias("h0SFD2k2#Fdsa"), Equals, false)
	c.Check(isValidAlias("0dslka-4"), Equals, false)
	c.Check(isValidAlias("-fdslka"), Equals, false)
}

func (s *TestSuite) TestHumanizedTime(c *C) {
	hTime := timeDurationToHumanizedTime(time.Duration(10) * time.Second)
	c.Assert(hTime.Minutes, Equals, int64(0))
	c.Assert(hTime.Hours, Equals, int64(0))
	c.Assert(hTime.Days, Equals, int64(0))

	hTime = timeDurationToHumanizedTime(time.Duration(10) * time.Minute)
	c.Assert(hTime.Hours, Equals, int64(0))
	c.Assert(hTime.Days, Equals, int64(0))

	hTime = timeDurationToHumanizedTime(time.Duration(10) * time.Hour)
	c.Assert(hTime.Days, Equals, int64(0))

	hTime = timeDurationToHumanizedTime(time.Duration(24) * time.Hour)
	c.Assert(hTime.Days, Not(Equals), int64(0))
}

func (s *TestSuite) TestVersions(c *C) {
	v1, e := version.NewVersion("1.7.1")
	v2, e := version.NewVersion("1.6.3")
	c.Assert(e, IsNil)
	constraint, e := version.NewConstraint(">= 1.7.1")
	c.Assert(e, IsNil)
	c.Assert(constraint.Check(v1), Equals, true)
	c.Assert(constraint.Check(v2), Equals, false)
}
