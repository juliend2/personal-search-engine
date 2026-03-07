# TODO


- [ ] Implement parsing for:
    - [ ] txt
    - [ ] word (docx)
    - [ ] excel (xlsx)
    - [ ] pdf
- [ ] Implement crawling of:
    - [ ] dropbox files
    - [ ] google drive files
    - [ ] notion documents
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


