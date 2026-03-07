README
======

# TODO

- [ ] Implement parsing for:
    - [x] txt
    - [ ] word -- with [godocx](https://github.com/gomutex/godocx)
    - [ ] excel -- with [excelize](https://github.com/qax-os/excelize)
    - [ ] pdf -- with [pdfcpu](https://github.com/pdfcpu/pdfcpu)
    - [ ] image (to index EXIF's `geopoint` at some point), for geo-targetted searches -- with [go-exiftool](https://github.com/barasher/go-exiftool)
- [ ] Implement crawling of:
    - [ ] dropbox files
    - [ ] google drive files
    - [ ] notion documents
    - [ ] (nicetohave) todoist
    - [ ] (nicetohave) google calendar
    - [ ] (nicetohave) gmail
- [ ] Implement an indexer that:
    - [ ] maps all the data sources (dropbox, notion, etc)
    - [ ] finds out the file type
    - [ ] execute the proper parsing strategy
- [ ] Implement a search engine
    - [ ] with a CLI interface
    - [ ] with a Web UI


# Crawling

## rclone

supporte au moins google drive et dropbox.

### pour voir les fichiers sur mon drive.

specifiquement le fichier `impots 2025`:
```sh
rclone ls remote: | grep impots
```

### pour sync:
```sh
rclone sync remote: /home/julien/GoogleDrive --progress --exclude "/Eglise/**"
```
(où "remote" est le nom du remote que j'ai pris pour mon google drive)

(ça gère la conversion de google doc vers docx automatiquement)


# Indexing

[bleve](https://github.com/blevesearch/bleve)

## id for the index

`<source>-<path>`

where `source` is a slug of the source, such as:

- `db` for dropbox
- `gd` for google drive
- `nt` for notion

and where `path` is either an URL or a (relative) file path to the resource.

# Searching

[bleve](https://github.com/blevesearch/bleve)

