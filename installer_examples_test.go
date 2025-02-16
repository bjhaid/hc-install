// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package install_test

import (
	"context"
	"log"

	install "github.com/bjhaid/hc-install"
	"github.com/bjhaid/hc-install/build"
	"github.com/bjhaid/hc-install/fs"
	"github.com/bjhaid/hc-install/product"
	"github.com/bjhaid/hc-install/releases"
	"github.com/bjhaid/hc-install/src"
	"github.com/hashicorp/go-version"
)

// Installation of a single exact version
func ExampleInstaller() {
	ctx := context.Background()
	i := install.NewInstaller()
	defer i.Remove(ctx)
	v1_3 := version.Must(version.NewVersion("1.3.7"))

	execPath, err := i.Install(ctx, []src.Installable{
		&releases.ExactVersion{
			Product: product.Terraform,
			Version: v1_3,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Terraform %s installed to %s", v1_3, execPath)

	// run any tests
}

// Locating or installing latest version per constraint
func ExampleInstaller_latestVersionConstrained() {
	ctx := context.Background()
	i := install.NewInstaller()
	defer i.Remove(ctx)

	v1 := version.MustConstraints(version.NewConstraint("~> 1.0"))

	execPath, err := i.Ensure(context.Background(), []src.Source{
		&fs.Version{
			Product:     product.Terraform,
			Constraints: v1,
		},
		&releases.LatestVersion{
			Product:     product.Terraform,
			Constraints: v1,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Terraform %s available at %s", v1, execPath)

	// run any tests
}

// Installation of multiple versions
func ExampleInstaller_multipleVersions() {
	ctx := context.Background()
	i := install.NewInstaller()
	defer i.Remove(ctx)

	v1_1 := version.Must(version.NewVersion("1.1.0"))
	execPath, err := i.Install(context.Background(), []src.Installable{
		&releases.ExactVersion{
			Product: product.Terraform,
			Version: v1_1,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Terraform %s available at %s", v1_1, execPath)

	// run any 1.1 tests

	v1_3 := version.Must(version.NewVersion("1.3.0"))
	execPath, err = i.Install(context.Background(), []src.Installable{
		&releases.ExactVersion{
			Product: product.Terraform,
			Version: v1_3,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Terraform %s available at %s", v1_3, execPath)

	// run any 1.3 tests
}

// Installation and building of multiple versions
func ExampleInstaller_installAndBuildMultipleVersions() {
	ctx := context.Background()
	i := install.NewInstaller()
	defer i.Remove(ctx)

	vc := version.MustConstraints(version.NewConstraint("~> 1.3"))
	rv := &releases.Versions{
		Product:     product.Terraform,
		Constraints: vc,
	}

	versions, err := rv.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	versions = append(versions, &build.GitRevision{
		Product: product.Terraform,
		Ref:     "HEAD",
	})

	for _, installableVersion := range versions {
		execPath, err := i.Ensure(context.Background(), []src.Source{
			installableVersion,
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Terraform %s installed to %s", installableVersion, execPath)
		// run any tests
	}
}

// Installation of a single exact enterprise version
func ExampleInstaller_enterpriseVersion() {
	ctx := context.Background()
	i := install.NewInstaller()
	defer i.Remove(ctx)
	v1_9 := version.Must(version.NewVersion("1.9.8"))
	licenseDir := "/some/path"

	execPath, err := i.Install(ctx, []src.Installable{
		&releases.ExactVersion{
			Product: product.Vault,
			Version: v1_9,
			Enterprise: &releases.EnterpriseOptions{ // specify that we want the enterprise version
				LicenseDir: licenseDir, // where license files should be placed (required for enterprise versions)
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Vault %s Enterprise installed to %s; license information installed to %s", v1_9, execPath, licenseDir)

	// run any tests
}
