package pkg

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

type version struct {
	major     int
	minor     int
	patch     int
	preString string
}

type Releaser struct {
	raw  string
	v    version
	next version
}

func (r *Releaser) HasPrerelease() bool {
	return len(r.v.preString) > 0
}

func (r *Releaser) String() (s string) {
	var b strings.Builder

	b.Write([]byte("v"))

	b.WriteString(strconv.Itoa(r.next.major))
	b.WriteString(".")
	b.WriteString(strconv.Itoa(r.next.minor))
	b.WriteString(".")
	b.WriteString(strconv.Itoa(r.next.patch))

	if len(r.next.preString) > 0 {
		b.Write([]byte(r.next.preString))
	}

	s = b.String()

	return
}

type ReleaserMod func(*Releaser)

func CreateNewReleaser(vString string, mods ...ReleaserMod) (r *Releaser, err error) {

	if !semver.IsValid(vString) {
		err = fmt.Errorf("invalid version string")
		return
	}

	r = &Releaser{}
	r.raw = vString

	simplifiedV := semver.Canonical(vString)

	vParts := strings.Split(simplifiedV, ".")

	preParts := strings.Split(vParts[2], "-")

	if len(preParts) > 1 {
		r.v.patch, err = strconv.Atoi(preParts[0])
	} else {
		r.v.patch, err = strconv.Atoi(vParts[2])
	}

	if err != nil {
		return
	}

	r.v.minor, err = strconv.Atoi(vParts[1])
	if err != nil {
		return
	}

	majorStr := strings.Replace(vParts[0], "v", "", -1)
	r.v.major, err = strconv.Atoi(majorStr)
	if err != nil {
		return
	}

	r.v.preString = semver.Prerelease(vString)

	r.next.major = r.v.major
	r.next.minor = r.v.minor
	r.next.patch = r.v.patch
	r.next.preString = r.v.preString

	for _, mod := range mods {
		mod(r)
	}

	return
}

func WithPrerelease(pre string) ReleaserMod {
	return func(r *Releaser) {
		r.next.preString = pre
	}
}

func WithPrereleaseIncrement() ReleaserMod {
	return func(r *Releaser) {
		preParts := strings.Split(r.v.preString, ".")
		prePointer, _ := strconv.Atoi(preParts[1])
		prePointer += 1
		preParts[1] = strconv.Itoa(prePointer)
		r.next.preString = strings.Join(preParts[:], ".")
	}
}

func WithMajorIncrement() ReleaserMod {
	return func(r *Releaser) {
		r.next.major += 1
	}
}

func WithMinorIncrement() ReleaserMod {
	return func(r *Releaser) {
		r.next.minor += 1
	}
}

func WithPatchIncrement() ReleaserMod {
	return func(r *Releaser) {
		r.next.patch += 1
	}
}

func WithMajorReset() ReleaserMod {
	return func(r *Releaser) {
		r.next.major = 0
	}
}

func WithMinorReset() ReleaserMod {
	return func(r *Releaser) {
		r.next.minor = 0
	}
}

func WithPatchReset() ReleaserMod {
	return func(r *Releaser) {
		r.next.patch = 0
	}
}

func WithPrereleaseReset() ReleaserMod {
	return func(r *Releaser) {
		preParts := strings.Split(r.v.preString, ".")
		// reset something like `beta.1` to `beta.0`
		preParts[1] = strconv.Itoa(0)
		r.next.preString = strings.Join(preParts[:], ".")
	}
}

func WithClearPrerelease() ReleaserMod {
	return func(r *Releaser) {
		r.next.preString = ""
	}
}
