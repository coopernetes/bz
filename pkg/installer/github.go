package installer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"compress/gzip"
	"archive/tar"

	"github.com/google/go-github/v56/github"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/compfile"
	"github.com/rwxrob/help"
)

var GithubCmd = &Z.Cmd{
	Name:    `github-installer`,
	Summary: `install a binary tool from GitHub releases`,
	Version: `v0.0.1`,
	License: `Apache-2.0`,
	Source:  `github.com/coopernetes/bz/pkg/installer`,
	Comp:    compfile.New(),
	Aliases: []string{`gh-install`, `gh-i`},
	UseVars: true,
	Params: []string{
		"version",
	},
	MinArgs: 1,
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(cmd *Z.Cmd, args ...string) error {
		version := cmd.Param("version")
		ghProject := strings.Split(args[0], "/")
		if len(ghProject) != 2 {
			return fmt.Errorf("invalid GitHub project name %v", ghProject)
		}
		owner := ghProject[0]
		repo := ghProject[1]
		if version == "" {
			version = "latest"
		}
		log.Printf("installing %s version %s", ghProject, version)
		ctx := context.Background()
		client := github.NewClient(nil)

		release, resp, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return errors.New(resp.Status)
		}
		log.Printf("found latest release %s", release.GetTagName())
		// 1. Get GOOS and GOARCH
		// 2. Get the asset name that contains GOOS and GOARCH
		// 3. Download the asset
		// 4. If a signature or checksum is available, validate it
		// 5. Install the asset to $HOME/.local/bin (unzip/untar or copy directly)
		goos, goarch := runtime.GOOS, runtime.GOARCH
		matchedAssets := make([]*github.ReleaseAsset, 0)
		for _, asset := range release.Assets {
			if strings.Contains(asset.GetName(), goos) && strings.Contains(asset.GetName(), goarch) {
				matchedAssets = append(matchedAssets, asset)
			}
		}
		if len(matchedAssets) == 0 {
			return fmt.Errorf("no binary found for %s/%s", goos, goarch)
		}
		log.Printf("matched assets: %v", matchedAssets)
		var gztar *github.ReleaseAsset
		var zip *github.ReleaseAsset
		for _, a := range matchedAssets {
			if strings.HasSuffix(a.GetName(), ".tar.gz") {
				gztar = a
			}
			if strings.HasSuffix(a.GetName(), ".zip") {
				zip = a
			}
		}
		if gztar == nil && zip == nil {
			return fmt.Errorf("no binary found for %s/%s", goos, goarch)
		}
		if gztar != nil {
			log.Printf("downloading %s", gztar.GetName())
			asset, _, err := client.Repositories.DownloadReleaseAsset(ctx, owner, repo, gztar.GetID(), http.DefaultClient)
			if err != nil {
				return err
			}
			defer asset.Close()
			file, err := os.CreateTemp("", "github-installer-")
			if err != nil {
				return err
			}
			defer file.Close()
			io.Copy(file, asset)
			file.Seek(0, 0)
			log.Printf("extracting %s", gztar.GetName())
			gzreader, err := gzip.NewReader(file)
			if err != nil {
				return err
			}
			defer gzreader.Close()
			tarreader := tar.NewReader(gzreader)

			curdir, err := os.Getwd()
			if err != nil {
				return err
			}
			outDir := curdir + "/.tmp"
			for {
				hdr, err := tarreader.Next()
				if err == io.EOF {
					break // End of archive
				}
				if err != nil {
					return err
				}
				log.Printf("tar file %s:\n", hdr.Name)
				if hdr.Typeflag == tar.TypeDir {
					if err := os.MkdirAll(outDir + "/" + hdr.Name, 0755); err != nil {
						return err
					}
				} else {
					n, err := os.Create(outDir + "/" + hdr.Name)
					if err != nil {
						return err
					}
					if _, err := io.Copy(n, tarreader); err != nil {
						return err
					}
				}
			}
			log.Printf("extracted %s to %s", gztar.GetName(), outDir)
		}

		return nil
	},
}

var GithubProjectVersion = &Z.Cmd{
	Name:    `github-project-version`,
	Summary: `finds the latest version of a GitHub project`,
	Call: func(cmd *Z.Cmd, args ...string) error {
		return nil
	},
}
