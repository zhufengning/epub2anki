package main

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/beevik/etree"
	_ "github.com/mattn/go-sqlite3"
)

var noUI *bool
var zipURI *string
var ankiURI *string
var db *sql.DB
var bookName string
var cnt = 0
var lb *widget.Label

func init() {
	noUI = flag.Bool("noui", false, "disable UI")
	zipURI = flag.String("i", "", "The path of your epub file")
	ankiURI = flag.String("o", "", "The path to save .apkg file")
}

func toAnki(zFile *zip.Reader) {
	f, e := zFile.Open("META-INF/container.xml")
	defer f.Close()
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		log.Fatal(e)
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	doc := etree.NewDocument()
	if e := doc.ReadFromBytes(buf.Bytes()); e != nil {
		lb.SetText(fmt.Sprint(e))
		log.Fatal(e)
		return
	}
	rootFiles := doc.FindElements("container/rootfiles/rootfile")
	for _, v := range rootFiles {
		for _, attr := range v.Attr {
			if attr.Key == "full-path" {
				readToc(zFile, attr.Value)
			}
		}
	}

}

func readToc(zFile *zip.Reader, path string) {
	toc, e := zFile.Open(path)
	defer toc.Close()
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		log.Fatal(e)
		return
	}
	doc := etree.NewDocument()
	body, e := io.ReadAll(toc)
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		log.Fatal(e)
		return
	}
	doc.ReadFromBytes(body)
	bookName = regexp.MustCompile(`<dc:title>.*?</dc:title>`).FindString(string(body))
	bookName = strings.Replace(bookName, "<dc:title>", "", 1)
	bookName = strings.Replace(bookName, "</dc:title>", "", 1)
	fmt.Println(bookName)
	db.Exec(fmt.Sprintf(`	
	INSERT INTO "main"."col" ("id", "crt", "mod", "scm", "ver", "dty", "usn", "ls", "conf", "models", "decks", "dconf", "tags") VALUES ('1', '1615060800', '1615304856121', '1615304855707', '11', '0', '0', '0', '{"activeDecks":[1],"curDeck":"1","newSpread":0,"collapseTime":1200,"timeLim":0,"estTimes":true,"dueCounts":true,"curModel":1615304855848,"nextPos":1,"sortType":"noteFld","sortBackwards":false,"addToCur":true,"newBury":true}', '{"1615304855845":{"sortf":0,"did":1,"latexPre":"\\documentclass[12pt]{article}\n\\special{papersize=3in,5in}\n\\usepackage[utf8]{inputenc}\n\\usepackage{amssymb,amsmath}\n\\pagestyle{empty}\n\\setlength{\\parindent}{0in}\n\\begin{document}\n","latexPost":"\\end{document}","mod":1615304855,"usn":-1,"vers":[],"type":0,"css":".card {\n font-family: arial;\n font-size: 20px;\n text-align: center;\n color: black;\n background-color: white;\n}","name":"基本型（输入答案）","flds":[{"name":"正面","ord":0,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]},{"name":"背面","ord":1,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]}],"tmpls":[{"name":"卡片 1","ord":0,"qfmt":"{{正面}}\n\n{{type:背面}}","afmt":"{{正面}}\n\n<hr id=answer>\n\n{{type:背面}}","did":null,"bqfmt":"","bafmt":"","bfont":"Arial","bsize":12}],"tags":[],"id":1615304855845,"req":[[0,"all",[0]]]},"1615304855841":{"sortf":0,"did":1,"latexPre":"\\documentclass[12pt]{article}\n\\special{papersize=3in,5in}\n\\usepackage[utf8]{inputenc}\n\\usepackage{amssymb,amsmath}\n\\pagestyle{empty}\n\\setlength{\\parindent}{0in}\n\\begin{document}\n","latexPost":"\\end{document}","mod":1615304855,"usn":-1,"vers":[],"type":0,"css":".card {\n font-family: arial;\n font-size: 20px;\n text-align: center;\n color: black;\n background-color: white;\n}","name":"基本型（含翻转的卡片）","flds":[{"name":"正面","ord":0,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]},{"name":"背面","ord":1,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]}],"tmpls":[{"name":"卡片 1","ord":0,"qfmt":"{{正面}}","afmt":"{{FrontSide}}\n\n<hr id=answer>\n\n{{背面}}","did":null,"bqfmt":"","bafmt":"","bfont":"Arial","bsize":12},{"name":"卡片 2","ord":1,"qfmt":"{{背面}}","afmt":"{{FrontSide}}\n\n<hr id=answer>\n\n{{正面}}","did":null,"bqfmt":"","bafmt":"","bfont":"Arial","bsize":12}],"tags":[],"id":1615304855841,"req":[[0,"all",[0]],[1,"all",[1]]]},"1615096069364":{"sortf":0,"did":1615096076200,"latexPre":"\\documentclass[12pt]{article}\n\\special{papersize=3in,5in}\n\\usepackage[utf8]{inputenc}\n\\usepackage{amssymb,amsmath}\n\\pagestyle{empty}\n\\setlength{\\parindent}{0in}\n\\begin{document}\n","latexPost":"\\end{document}","mod":1615096069,"usn":-1,"vers":[],"type":0,"css":".card {\n font-family: arial;\n font-size: 20px;\n text-align: center;\n color: black;\n background-color: white;\n}","name":"基本","flds":[{"name":"正面","ord":0,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]},{"name":"背面","ord":1,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]}],"tmpls":[{"name":"卡片 1","ord":0,"qfmt":"{{正面}}","afmt":"{{FrontSide}}\n\n<hr id=answer>\n\n{{背面}}","did":null,"bqfmt":"","bafmt":"","bfont":"Arial","bsize":12}],"tags":[],"id":1615096069364,"req":[[0,"all",[0]]]},"1615304855832":{"sortf":0,"did":1,"latexPre":"\\documentclass[12pt]{article}\n\\special{papersize=3in,5in}\n\\usepackage[utf8]{inputenc}\n\\usepackage{amssymb,amsmath}\n\\pagestyle{empty}\n\\setlength{\\parindent}{0in}\n\\begin{document}\n","latexPost":"\\end{document}","mod":1615304855,"usn":-1,"vers":[],"type":1,"css":".card {\n font-family: arial;\n font-size: 20px;\n text-align: center;\n color: black;\n background-color: white;\n}.cloze {font-weight: bold;color: blue;}","name":"完形填空","flds":[{"name":"文本","ord":0,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]},{"name":"更多","ord":1,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]}],"tmpls":[{"name":"完形填空","ord":0,"qfmt":"{{cloze:文本}}","afmt":"{{cloze:文本}}<br>\n{{更多}}","did":null,"bqfmt":"","bafmt":"","bfont":"Arial","bsize":12}],"tags":[],"id":1615304855832},"1615304855848":{"sortf":0,"did":1,"latexPre":"\\documentclass[12pt]{article}\n\\special{papersize=3in,5in}\n\\usepackage[utf8]{inputenc}\n\\usepackage{amssymb,amsmath}\n\\pagestyle{empty}\n\\setlength{\\parindent}{0in}\n\\begin{document}\n","latexPost":"\\end{document}","mod":1615304855,"usn":-1,"vers":[],"type":0,"css":".card {\n font-family: arial;\n font-size: 20px;\n text-align: center;\n color: black;\n background-color: white;\n}","name":"基本","flds":[{"name":"正面","ord":0,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]},{"name":"背面","ord":1,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]}],"tmpls":[{"name":"卡片 1","ord":0,"qfmt":"{{正面}}","afmt":"{{FrontSide}}\n\n<hr id=answer>\n\n{{背面}}","did":null,"bqfmt":"","bafmt":"","bfont":"Arial","bsize":12}],"tags":[],"id":1615304855848,"req":[[0,"all",[0]]]},"1615304855833":{"sortf":0,"did":1,"latexPre":"\\documentclass[12pt]{article}\n\\special{papersize=3in,5in}\n\\usepackage[utf8]{inputenc}\n\\usepackage{amssymb,amsmath}\n\\pagestyle{empty}\n\\setlength{\\parindent}{0in}\n\\begin{document}\n","latexPost":"\\end{document}","mod":1615304855,"usn":-1,"vers":[],"type":0,"css":".card {\n font-family: arial;\n font-size: 20px;\n text-align: center;\n color: black;\n background-color: white;\n}","name":"基本型（随意添加翻转的卡片）","flds":[{"name":"正面","ord":0,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]},{"name":"背面","ord":1,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]},{"name":"添加翻转卡片","ord":2,"sticky":false,"rtl":false,"font":"Arial","size":20,"media":[]}],"tmpls":[{"name":"卡片 1","ord":0,"qfmt":"{{正面}}","afmt":"{{FrontSide}}\n\n<hr id=answer>\n\n{{背面}}","did":null,"bqfmt":"","bafmt":"","bfont":"Arial","bsize":12},{"name":"卡片 2","ord":1,"qfmt":"{{#添加翻转卡片}}{{背面}}{{/添加翻转卡片}}","afmt":"{{FrontSide}}\n\n<hr id=answer>\n\n{{正面}}","did":null,"bqfmt":"","bafmt":"","bfont":"Arial","bsize":12}],"tags":[],"id":1615304855833,"req":[[0,"all",[0]],[1,"all",[1,2]]]}}', '{"1":{"newToday":[0,0],"revToday":[0,0],"lrnToday":[0,0],"timeToday":[0,0],"conf":1,"usn":0,"desc":"","dyn":0,"collapsed":false,"extendNew":10,"extendRev":50,"id":1,"name":"Default","mod":1615304855},"1615096076200":{"newToday":[2,0],"revToday":[2,0],"lrnToday":[2,0],"timeToday":[2,0],"conf":1,"usn":-1,"desc":"","dyn":0,"collapsed":false,"extendNew":10,"extendRev":50,"name":%q,"id":1615096076200,"mod":1615096076}}', '{"1":{"name":"Default","new":{"delays":[1,10],"ints":[1,4,7],"initialFactor":2500,"separate":true,"order":1,"perDay":20,"bury":false},"lapse":{"delays":[10],"mult":0,"minInt":1,"leechFails":8,"leechAction":0},"rev":{"perDay":100,"ease4":1.3,"fuzz":0.05,"minSpace":1,"ivlFct":1,"maxIvl":36500,"bury":false},"maxTaken":60,"timer":0,"autoplay":true,"replayq":true,"mod":0,"usn":0,"id":1}}', '{}');
	`, bookName))
	itemMap := make(map[string]string)
	for _, item := range doc.FindElements("package/manifest/item") {
		var id, href string
		for _, attr := range item.Attr {
			if attr.Key == "id" {
				id = attr.Value
			} else if attr.Key == "href" {
				href = attr.Value
			}
		}
		itemMap[id] = href
	}
	for _, chap := range doc.FindElements("package/spine/itemref") {
		for _, idref := range chap.Attr {
			if idref.Key == "idref" {
				cnt += 1
				addToDB(zFile, itemMap[idref.Value])
			}
		}
	}
	cnt = 0
}

func addToDB(zFile *zip.Reader, path string) {
	f, e := zFile.Open(path)
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		log.Println(e)
		path = "OEBPS/" + path
		f, e = zFile.Open(path)
	}
	defer f.Close()
	nr, e := io.ReadAll(f)
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		log.Fatal(e)
		return
	}
	doc := string(nr)
	title := regexp.MustCompile(`<title>.*?</title>`).FindString(doc)
	title = strings.Replace(title, "<title>", "", 1)
	title = strings.Replace(title, "/<title>", "", 1)
	title = fmt.Sprint(cnt) + title

	t := time.Now().UnixNano() / 1000000
	s := fmt.Sprintf(
		`
	INSERT INTO "main"."cards" ("id", "nid", "did", "ord", "mod", "usn", "type", "queue", "due", "ivl", "factor", "reps", "lapses", "left", "odue", "odid", "flags", "data") VALUES ('%d', '%d', '1615096076200', '0', '%d', '-1', '0', '0', '1', '0', '0', '0', '0', '0', '0', '0', '0', '');
	`, t, t, t/1000)
	_, e = db.Exec(s)
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		log.Fatal(e)
	}
	doc = strings.ReplaceAll(doc, "\"", "'")
	flds := title + string(rune(31)) + doc
	flds = strings.ReplaceAll(flds, "\"", "'")

	s = fmt.Sprintf(`
	INSERT INTO "main"."notes" ("id", "guid", "mid", "mod", "usn", "tags", "flds", "sfld", "csum", "flags", "data") VALUES ('%d', '%d', '1615096069364', '%d', '-1', '', "%s", "%s", '%d', '0', '');
	`, t, t, t/1000, flds, doc, 1234)
	fmt.Println(s)
	_, e = db.Exec(s)
	if e != nil {
		log.Fatal(e)
	}
	time.Sleep(1 * time.Millisecond)
}

func packApkg() {
	path := os.TempDir()
	f, e := os.Create(path + "/out.apkg")
	defer f.Close()
	if e != nil {
		lb.SetText(fmt.Sprint(e))
	}
	z := zip.NewWriter(f)
	m, e := z.Create("media")
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		return
	}
	_, e = io.WriteString(m, "{}")
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		return
	}
	c, e := z.Create("collection.anki2")
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		return
	}
	anki, e := os.Open(path + "/collection.anki2")
	defer anki.Close()
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		return
	}
	nr, e := io.ReadAll(anki)
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		return
	}
	_, e = c.Write(nr)
	if e != nil {
		lb.SetText(fmt.Sprint(e))
		return
	}
	z.Close()
}

func main() {
	path := os.TempDir()
	flag.Parse()
	if *noUI {
		initDB()
		z, e := os.OpenFile(*zipURI, os.O_RDONLY, 0666)
		if e != nil {
			log.Fatal(e)
		}
		zBytes, e := io.ReadAll(z)
		zReader := bytes.NewReader(zBytes)
		zFile, e := zip.NewReader(zReader, int64(len(zBytes)))
		toAnki(zFile)
		db.Close()
		packApkg()
		out, e := os.Create(*ankiURI)
		defer out.Close()
		if e != nil {
			log.Fatal(e)
		}
		res, e := os.Open(path + "/out.apkg")
		defer res.Close()
		if e != nil {
			log.Fatal(e)
		}
		b, e := io.ReadAll(res)
		if e != nil {
			log.Fatal(e)
		}
		_, e = out.Write(b)
		if e != nil {
			log.Fatal(e)
		}
	} else {
		a := app.New()
		w := a.NewWindow("Hello")
		//a.Settings().SetTheme(&Biu{})
		//uncomment to build with theme(chinese font)
		but := widget.NewButton("Choose File", func() {})
		lb = widget.NewLabel("waiting")
		but.OnTapped = func() {
			dg := dialog.NewFileOpen(func(uri fyne.URIReadCloser, e error) {
				lb.SetText("Running")
				but.Disable()
				if uri == nil {
					lb.SetText("no file choosed")
					but.Enable()
					return
				}
				initDB()
				if e != nil {
					lb.SetText(fmt.Sprint(e))
					but.Enable()
					return
				}
				b, e := io.ReadAll(uri)
				zReader := bytes.NewReader(b)
				zFile, e := zip.NewReader(zReader, int64(len(b)))
				toAnki(zFile)
				db.Close()
				dialog.NewFileSave(func(uri fyne.URIWriteCloser, e error) {
					packApkg()
					res, e := os.Open(path + "/out.apkg")
					defer res.Close()
					if e != nil {
						lb.SetText(fmt.Sprint(e))
						but.Enable()
						return
					}
					b, e := io.ReadAll(res)
					if e != nil {
						lb.SetText(fmt.Sprint(e))
						but.Enable()
						return
					}
					_, e = uri.Write(b)
					if e != nil {
						lb.SetText(fmt.Sprint(e))
						but.Enable()
						return
					}
					uri.Close()
				}, w).Show()
				but.Enable()
				lb.SetText("finished")
			}, w)
			dg.Show()
		}
		w.SetContent(container.NewScroll(container.NewVBox(but, lb)))
		w.Resize(fyne.NewSize(640.0, 360.0))
		w.ShowAndRun()
	}
}

func initDB() {
	path := os.TempDir()
	var e error
	e = os.Remove(path + "/collection.anki2")
	// if e != nil {
	// 	lb.SetText(fmt.Sprint(e))
	// }
	db, e = sql.Open("sqlite3", os.TempDir()+"/collection.anki2")
	if e != nil {
		lb.SetText(fmt.Sprint(e))
	}

	//t := time.Now().UnixNano() / 1000000
	s := fmt.Sprintf(
		`
		CREATE TABLE revlog (
			id              integer primary key,
			   -- epoch-milliseconds timestamp of when you did the review
			cid             integer not null,
			   -- cards.id
			usn             integer not null,
				-- update sequence number: for finding diffs when syncing. 
				--   See the description in the cards table for more info
			ease            integer not null,
			   -- which button you pushed to score your recall. 
			   -- review:  1(wrong), 2(hard), 3(ok), 4(easy)
			   -- learn/relearn:   1(wrong), 2(ok), 3(easy)
			ivl             integer not null,
			   -- interval (i.e. as in the card table)
			lastIvl         integer not null,
			   -- last interval (i.e. the last value of ivl. Note that this value is not necessarily equal to the actual interval between this review and the preceding review)
			factor          integer not null,
			  -- factor
			time            integer not null,
			   -- how many milliseconds your review took, up to 60000 (60s)
			type            integer not null
			   --  0=learn, 1=review, 2=relearn, 3=cram
		);
		
		CREATE TABLE graves (
			usn             integer not null,
			oid             integer not null,
			type            integer not null
		);
		
	CREATE TABLE cards (
		id              integer primary key,
			-- the epoch milliseconds of when the card was created
		nid             integer not null,--    
			-- notes.id
		did             integer not null,
			-- deck id (available in col table)
		ord             integer not null,
			-- ordinal : identifies which of the card templates or cloze deletions it corresponds to 
			--   for card templates, valid values are from 0 to num templates - 1
		mod             integer not null,
			-- modificaton time as epoch seconds
		usn             integer not null,
			-- update sequence number : used to figure out diffs when syncing. 
			--   value of -1 indicates changes that need to be pushed to server. 
			--   usn < server usn indicates changes that need to be pulled from server.
		type            integer not null,
			-- 0=new, 1=learning, 2=review, 3=relearning
		queue           integer not null,
			-- -3=user buried(In scheduler 2),
			-- -2=sched buried (In scheduler 2), 
			-- -2=buried(In scheduler 1),
			-- -1=suspended,
			-- 0=new, 1=learning, 2=review (as for type)
			-- 3=in learning, next rev in at least a day after the previous review
			-- 4=preview
		due             integer not null,
			-- Due is used differently for different card types: 
			--   new: note id or random int
			--   due: integer day, relative to the collection's creation time
			--   learning: integer timestamp in second
		ivl             integer not null,
			-- interval (used in SRS algorithm). Negative = seconds, positive = days
		factor          integer not null,
			-- The ease factor of the card in permille (parts per thousand). If the ease factor is 2500, the card’s interval will be multiplied by 2.5 the next time you press Good.
		reps            integer not null,
			-- number of reviews
		lapses          integer not null,
			-- the number of times the card went from a "was answered correctly" 
			--   to "was answered incorrectly" state
		left            integer not null,
			-- of the form a*1000+b, with:
			-- b the number of reps left till graduation
			-- a the number of reps left today
		odue            integer not null,
			-- original due: In filtered decks, it's the original due date that the card had before moving to filtered.
						-- If the card lapsed in scheduler1, then it's the value before the lapse. (This is used when switching to scheduler 2. At this time, cards in learning becomes due again, with their previous due date)
						-- In any other case it's 0.
		odid            integer not null,
			-- original did: only used when the card is currently in filtered deck
		flags           integer not null,
			-- an integer. This integer mod 8 represents a "flag", which can be see in browser and while reviewing a note. Red 1, Orange 2, Green 3, Blue 4, no flag: 0. This integer divided by 8 represents currently nothing
		data            text not null
			-- currently unused
	);
	CREATE TABLE IF NOT EXISTS "col" (
		"id"	integer,
		"crt"	integer NOT NULL,
		"mod"	integer NOT NULL,
		"scm"	integer NOT NULL,
		"ver"	integer NOT NULL,
		"dty"	integer NOT NULL,
		"usn"	integer NOT NULL,
		"ls"	integer NOT NULL,
		"conf"	text NOT NULL,
		"models"	text NOT NULL,
		"decks"	text NOT NULL,
		"dconf"	text NOT NULL,
		"tags"	text NOT NULL,
		PRIMARY KEY("id")
	);
	CREATE TABLE notes (
		id              integer primary key,
		  -- epoch miliseconds of when the note was created
		guid            text not null,
		  -- globally unique id, almost certainly used for syncing
		mid             integer not null,
		  -- model id
		mod             integer not null,
		  -- modification timestamp, epoch seconds
		usn             integer not null,
		  -- update sequence number: for finding diffs when syncing.
		  --   See the description in the cards table for more info
		tags            text not null,
		  -- space-separated string of tags. 
		flds            text not null,
		  -- the values of the fields in this note. separated by 0x1f (31) character.
		sfld            integer not null,
		  -- sort field: used for quick sorting and duplicate check. The sort field is an integer so that when users are sorting on a field that contains only numbers, they are sorted in numeric instead of lexical order. Text is stored in this integer field.
		csum            integer not null,
		  -- field checksum used for duplicate check.
		  --   integer representation of first 8 digits of sha1 hash of the first field
		flags           integer not null,
		  -- unused
		data            text not null
		  -- unused
	);
	CREATE INDEX ix_cards_nid on cards (nid);
	CREATE INDEX ix_cards_sched on cards (did, queue, due);
	CREATE INDEX ix_cards_usn on cards (usn);
	CREATE INDEX ix_notes_csum on notes (csum);
	CREATE INDEX ix_notes_usn on notes (usn);
	CREATE INDEX ix_revlog_cid on revlog (cid);
	CREATE INDEX ix_revlog_usn on revlog (usn);
	`)

	_, e = db.Exec(s)

	if e != nil {
		lb.SetText(fmt.Sprint(e))
	}
}
