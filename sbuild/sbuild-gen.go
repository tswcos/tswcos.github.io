/*
 Copyright (C) 2017 Toshiba Corporation

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 2 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful, but
 WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see
 <http://www.gnu.org/licenses/>.
*/

package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	sbuildLog            = "https://raw.githubusercontent.com/tswcos/build-log/master/"
	debian_cross_patches = "https://github.com/meta-debian/debian-cross-patches/tree/master/"
	white                = "#ecf0f1"
	green                = "#2ecc71"
	orange               = "#e67e22"
	blue                 = "#86c1b9"
)

func genDCB(notes map[string]interface{}) {
	pkgListFile, _ := ioutil.ReadFile("debian-cross-patches-list.txt")
	lines := strings.Split(string(pkgListFile), "\n")

	html := ""
	html += "<html><head><title>debian-cross-patches</title>\n"
	html += "<script src=\"https://kryogenix.org/code/browser/sorttable/sorttable.js\"></script>\n"
	html += "<script src=\"../js/default.js\"></script>\n"
	html += "<link rel=\"stylesheet\" type=\"text/css\" id=\"table_row_counter\" href=\"../css/table_row_counter.css\"/>\n"
	html += "</head><body>\n"
	html += "<h1>debian-cross-patches summary</h1>\n"
	html += "<h3><a href=\"https://github.com/meta-debian/debian-cross-patches\">Repository</a></h3>"
	html += "<b>Submitted</b> : patch has been submitted by meta-debian team.<br/>"
	html += "<b>Accepted</b>  : patch has been accepted by Debian.<br/>"
	html += "<br/><i>Packages, which have bug report but have not been submitted, mean they have been fixed by others.</i>"
	html += "<br/>"
	html += "<br/><table class=\"sortable\" id=\"sortable\">\n"
	html += "<tr bgcolor=\"#bdc3c7\">" +
		"<th class=\"sorttable_nosort\"></th>" +
		"<th>Source Name</th>" +
		"<th>Submitted</th>" +
		"<th>Accepted</th>" +
		"<th>Debian Bug</th>" +
		"<th>Remark</th></tr>\n"

	for i := 0; i < len(lines)-1; i++ {
		pkg := lines[i]

		remark := ""
		bugs := ""
		submitted := ""
		accepted := ""
		bgcolor := white
		if notes[pkg] != nil {
			data := notes[pkg].(map[string]interface{})
			if data["remark"] != nil {
				remark = data["remark"].(string)
			}

			if data["bugs"] != nil {
				bugsArr := data["bugs"].([]interface{})
				if len(bugsArr) > 0 {
					bgcolor = orange
				}
				for j := 0; j < len(bugsArr); j++ {
					bugs += "<a href=\"https://bugs.debian.org/" + bugsArr[j].(string) +
						"\">#" + bugsArr[j].(string) + "</a>"
					if j != len(bugsArr)-1 {
						bugs += "<br/>"
					}
				}
			}

			if data["submitted"] != nil {
				if data["submitted"].(bool) {
					submitted = "Yes"
					bgcolor = blue
				}
				if data["accepted"] != nil {
					if data["accepted"].(bool) {
						accepted = "Yes"
						bgcolor = green
					}
				}
			}

		}

		html += "<tr bgcolor=\"" + bgcolor + "\"><td></td>\n" +
			"<td>" + pkg + "</td>\n" +
			"<td>" + submitted + "</td>\n" +
			"<td>" + accepted + "</td>\n" +
			"<td>" + bugs + "</td>\n" +
			"<td>" + remark + "</td>\n" +
			"</tr>"
	}

	html += "</table></body></html>"

	f, _ := os.Create("debian-cross-patches.html")
	defer f.Close()
	f.WriteString(html)
	f.Sync()
	w := bufio.NewWriter(f)
	w.Flush()
}

func genIndex(notes map[string]interface{}) {
	statusFile, _ := ioutil.ReadFile("sbuild-result")
	lines := strings.Split(string(statusFile), "\n")

	html := "<html><head><title>Sbuild Status</title>\n"
	html += "<script src=\"https://kryogenix.org/code/browser/sorttable/sorttable.js\"></script>\n"
	html += "<script src=\"../js/default.js\"></script>\n"
	html += "<link rel=\"stylesheet\" type=\"text/css\" id=\"table_row_counter\" href=\"../css/table_row_counter.css\"/>\n"
	html += "</head><body>\n"
	html += "<h1>Debian cross-build state</h1>\n"
	html += "<h3><a href=\"./debian-cross-patches.html\">Summary of debian-cross-patches</a></h3>\n"
	html += "Build Architecture: amd64<br/>Host Architecture: armhf<br/>##Summary##<br/>\n"
	html += "<br/><input type=\"checkbox\"/ onclick=\"tableRemoveCounter(this);\"> Disable row counter. This helps table sort faster. (Re-enabling will take time)\n"
	html += "<br/><table class=\"sortable\" id=\"sortable\">\n"
	html += "<tr bgcolor=\"#bdc3c7\">" +
		"<th class=\"sorttable_nosort\"></th>" +
		"<th>Source Name</th>" +
		"<th class=\"sorttable_nosort\">Version</th>" +
		"<th>Status</th>" +
		"<th>Build At</th>" +
		"<th>Debian Bug</th>" +
		"<th>Remark</th></tr>\n"

	successCount := 0

	for i := 0; i < len(lines)-1; i++ {
		pkgInfo := strings.Fields(lines[i])
		name := pkgInfo[0]
		version := pkgInfo[1]
		status := pkgInfo[2]
		t, _ := time.Parse(time.RFC3339, pkgInfo[3])
		timestmp := t.Format("Jan _2, 2006")
		customTimestmp := t.Format("20060102150405")

		bgcolor := "#df2029"
		if status == "attempted" {
			bgcolor = "#e74c3c"
		} else if status == "skipped" {
			bgcolor = white
		} else if status == "successful" {
			bgcolor = green
			successCount++
		} else if status == "given-back" {
			bgcolor = orange
		}

		remark := ""
		bugs := ""
		if notes[name] != nil {
			data := notes[name].(map[string]interface{})
			if data["remark"] != nil {
				remark = data["remark"].(string)
			}
			if data["bugs"] != nil {
				bugsArr := data["bugs"].([]interface{})
				for j := 0; j < len(bugsArr); j++ {
					bugs += "<a href=\"https://bugs.debian.org/" + bugsArr[j].(string) +
						"\">#" + bugsArr[j].(string) + "</a>"
					if j != len(bugsArr)-1 {
						bugs += "<br/>"
					}
				}
			}
		}
		if _, err := os.Stat("./debian-cross-patches/" + name); !os.IsNotExist(err) {
			if remark != "" {
				remark += "<br/>"
			}
			remark += "\n (<a href=\"" + debian_cross_patches + name + "\">debian-cross-patches</a>)"
		}

		logFile := name + "_" + version + "_armhf.build"
		row := "<tr bgcolor=\"" + bgcolor + "\"><td></td>\n" +
			"<td><a href=\"" + sbuildLog + logFile + "\">" + name + "</a></td>\n" +
			"<td>" + version + "</td>\n" +
			"<td>" + status + "</td>\n" +
			"<td sorttable_customkey=\"" + customTimestmp + "\">" + timestmp + "</td>\n" +
			"<td>" + bugs + "</td>\n" +
			"<td>" + remark + "</td></tr>\n"
		html += row
	}

	html = strings.Replace(html, "##Summary##", "Success: "+strconv.Itoa(successCount)+"<br/>Total: "+strconv.Itoa(len(lines)-1), -1)
	html += "</table></body></html>"

	f, _ := os.Create("index.html")
	defer f.Close()
	f.WriteString(html)
	f.Sync()
	w := bufio.NewWriter(f)
	w.Flush()
}

func main() {
	var notes map[string]interface{}
	notesFile, _ := ioutil.ReadFile("notes.json")
	json.Unmarshal(notesFile, &notes)

	genIndex(notes)
	genDCB(notes)
}
