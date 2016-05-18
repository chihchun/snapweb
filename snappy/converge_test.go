/*
 * Copyright (C) 2014-2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package snappy

import (
	"errors"
	"os"

	"github.com/ubuntu-core/snappy/client"
	"github.com/ubuntu-core/snappy/snap"
	. "gopkg.in/check.v1"
	"launchpad.net/webdm/statustracker"
)

type PackagePayloadSuite struct {
	h Handler
	c *fakeSnapdClient
}

var _ = Suite(&PackagePayloadSuite{})

func (s *PackagePayloadSuite) SetUpTest(c *C) {
	os.Setenv("SNAP_DATA", c.MkDir())
	s.h.statusTracker = statustracker.New()
	s.c = &fakeSnapdClient{}
	s.h.setClient(s.c)
}

func (s *PackagePayloadSuite) TestPackageNotFound(c *C) {
	s.c.err = errors.New("the snap could not be retrieved")

	_, err := s.h.packagePayload("chatroom")
	c.Assert(err, NotNil)
}

func (s *PackagePayloadSuite) TestPackage(c *C) {
	s.c.snaps = []*client.Snap{newDefaultSnap()}

	pkg, err := s.h.packagePayload("chatroom")
	c.Assert(err, IsNil)
	c.Assert(pkg, DeepEquals, snapPkg{
		ID:            "chatroom",
		Description:   "WebRTC Video chat server for Snappy",
		DownloadSize:  0,
		Icon:          "/icons/chatroom_icon.png",
		InstalledSize: 18976651,
		Name:          "chatroom",
		Developer:     "ogra",
		Status:        "installed",
		Type:          "app",
		Version:       "0.1-8",
	})
}

type PayloadSuite struct {
	h Handler
}

var _ = Suite(&PayloadSuite{})

func (s *PayloadSuite) SetUpTest(c *C) {
	os.Setenv("SNAP_DATA", c.MkDir())
	s.h.statusTracker = statustracker.New()
	s.h.setClient(&fakeSnapdClient{})
}

func (s *PayloadSuite) TestPayload(c *C) {
	fakeSnap := newDefaultSnap()

	q := s.h.snapToPayload(fakeSnap)

	c.Check(q.Name, Equals, fakeSnap.Name)
	c.Check(q.Version, Equals, fakeSnap.Version)
	c.Check(q.Status, Equals, statustracker.StatusInstalled)
	c.Check(q.Type, Equals, snap.Type(fakeSnap.Type))
	c.Check(q.Icon, Equals, "/icons/chatroom_icon.png")
	c.Check(q.Description, Equals, fakeSnap.Description)
}

func (s *PayloadSuite) TestPayloadSnapInstalling(c *C) {
	fakeSnap := newDefaultSnap()
	fakeSnap.Status = client.StatusAvailable
	s.h.statusTracker.TrackInstall(fakeSnap)

	payload := s.h.snapToPayload(fakeSnap)
	c.Assert(payload.Status, Equals, statustracker.StatusInstalling)
}

type AllPackagesSuite struct {
	c *fakeSnapdClient
	h Handler
}

var _ = Suite(&AllPackagesSuite{})

func (s *AllPackagesSuite) SetUpTest(c *C) {
	os.Setenv("SNAP_DATA", c.MkDir())
	s.h.statusTracker = statustracker.New()
	s.c = &fakeSnapdClient{}
	s.h.setClient(s.c)
}

func (s *AllPackagesSuite) TestNoSnaps(c *C) {
	s.c.err = errors.New("snaps could not be filtered")

	snaps, err := s.h.allPackages(true, "")
	c.Assert(snaps, IsNil)
	c.Assert(err, NotNil)
}

func (s *AllPackagesSuite) TestHasSnaps(c *C) {
	s.c.snaps = []*client.Snap{
		newSnap("app2"),
		newSnap("app1"),
	}

	snaps, err := s.h.allPackages(true, "")
	c.Assert(err, IsNil)
	c.Assert(snaps, HasLen, 2)
	c.Assert(snaps[0].Name, Equals, "app1")
	c.Assert(snaps[1].Name, Equals, "app2")
}
