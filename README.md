# Go based Filesorting Tool

The idea behind this program is to sort PDF (and later maybe oder file types) with self defined tags. These tags are matched mostly with regular expressions against the content, filename, meta information and other factors that will depend on der file type.

This is my first project writen in GoLang and i'm realy looking for to any feedback that helps me to writ better code.

I was inspired by the "synOCR" and wanted to wirte somthing that can be used with a Synology or other simple file structure and does not have the tesseract as a core part of the program.



## implemented Features

- reading any configruation with a clean yaml file
  - Generate config if it does not exist
  - some parts can grow without problem like the list of tags that will be matched
- logging to stdout
- A simple FileQueue that handls the files that are to be processed
  - some functions that handls adding and removing files
- importing files in to the Queue
  - from a single source directory
- detecting filetype by contentType and filename suffix
- reading pdf content

## Comming Features

- logging to a logfile
- detecting tags that should be matched
- implementing some Features that handle the files based on the tags
- docker container
  - maybe with with and without teseract inside
- []features


## sources and links

- <https://www.synology-forum.de/threads/synocr-gui-fuer-ocrmypdf.99647/>