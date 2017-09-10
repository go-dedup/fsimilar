
# fsimilar

[![MIT License](http://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/go-dedup/fsimilar?status.svg)](http://godoc.org/github.com/go-dedup/fsimilar)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-dedup/fsimilar)](https://goreportcard.com/report/github.com/go-dedup/fsimilar)
[![travis Status](https://travis-ci.org/go-dedup/fsimilar.svg?branch=master)](https://travis-ci.org/go-dedup/fsimilar)

## TOC
- [fsimilar - find/file similar](#fsimilar---findfile-similar)
- [Synopsis](#synopsis)
- [Prefix, the near-duplicate detection algorithms](#prefix-the-near-duplicate-detection-algorithms)
  - [How to Identify Duplication](#how-to-identify-duplication)
  - [Near-duplicate detection algorithms](#near-duplicate-detection-algorithms)
    - [go-dedup collection](#go-dedup-collection)
    - [MinHash](#minhash)
    - [SimHash](#simhash)
    - [Vector space model](#vector-space-model)
    - [Advantages of fsimilar over existing tools](#advantages-of-fsimilar-over-existing-tools)
    - [Phonetic algorithms](#phonetic-algorithms)
  - [Using fsimilar](#using-fsimilar)
    - [Online usage help](#online-usage-help)
      - [$ fsimilar](#-fsimilar)
      - [$ fsimilar sim](#-fsimilar-sim)
      - [$ fsimilar vec](#-fsimilar-vec)
    - [Usage example](#usage-example)
      - [> test/sim.lstA](#-testsimlsta)
      - [$ fsimilar sim -i test/sim.lstA -d 12](#-fsimilar-sim--i-testsimlsta--d-12)
      - [$ fsimilar vec -i test/sim.lstA -t 0.7](#-fsimilar-vec--i-testsimlsta--t-07)
      - [$ fsimilar vec -i test/sim.lstA -t 0.7 -p](#-fsimilar-vec--i-testsimlsta--t-07--p)
      - [alternative invocation](#alternative-invocation)
    - [Using `fsimilar`](#using-`fsimilar`)
      - [> head -58 test/shell_test.ref](#-head--58-testshell_testref)
    - [The recommendations and decisions](#the-recommendations-and-decisions)
      - [> head -18 test/shell_ln.tmpl.sh](#-head--18-testshell_lntmplsh)
      - [> head -18 test/shell_mv.tmpl.sh](#-head--18-testshell_mvtmplsh)
      - [The approach](#the-approach)
    - [Binary releases](#binary-releases)
- [Download/Install](#downloadinstall)
  - [Download binaries](#download-binaries)
  - [Debian package](#debian-package)
  - [Install Source](#install-source)
- [Similar Projects](#similar-projects)
- [Credits](#credits)
- [Author(s) & Contributor(s)](#author(s)-&-contributor(s))

## fsimilar - find/file similar

`fsimilar` find similar file very fast -- it can spot similar files within the file system very quickly.

It is a useful tool to keep your file system slim.

# Synopsis

Take a look at [this file-list](https://github.com/go-dedup/fsimilar/blob/master/test/test1.lst), the first column is the file size and the next is the file specification itself, i.e., `/path/filename.ext`. See how easy or difficult for you to identify which files are actually the same. The above list is collected from a _single_ source, and there are already many duplicate files there. If you have ever collected different files from different sources, you will know that duplicate files are inevitable in our file system.

Identifying same files are easy for computers. However, identifying similar ones, i.e., the Near-Duplicate files turns out to be a hard problem even for computers, because it is no longer a zero-and-one, blank-and-white problem any more. The computer need to make guessing/decisions to get _rough estimates_. This is where it is hard for computers. In other words, the Near-Duplicate detection is in the realm of AI technology.

This project has two purposes,

- The tool solves the near-duplicate files problem by using the state of the art existing cutting-edge technologies.
- The source showcases the Go based near-duplicate detection algorithms collected in the [go-dedup](https://github.com/go-dedup) ecosystem.

Before we look at _what_ the tools is doing, let’s take a look at _how_ it is doing it first.

# Prefix, the near-duplicate detection algorithms

There is a famous article by mozilla on the [near-duplicate detection problem and algorithms](https://moz.com/devblog/near-duplicate-detection/). Let me quote the algorithm listing part as the quick summary, but please do give the article a careful read as it contains the most valuable explanation for the algorithms and their evolvements.

## How to Identify Duplication

There are many different ways that machines (that is, search engines and Moz) can attempt to identify duplicate content.

- **Bag of Words** – Comparing the words and frequency of those words on a page with those on another page.
- **Shingling** – This method improves on the “Bag of Words” approach by comparing short phrases, providing some context to the words.
- **Hashing** – This method further improves the process by eliminating the need to store copies of all of the content. The phrases are hashed into numbers, which can then be compared to identify duplication.
- **MinHash** – Helps streamline the process of storing content hashes.
- **SimHash** – Further streamlines the process of storing content hashes and duplication detection.

## Near-duplicate detection algorithms

### go-dedup collection

All the above listed algorithms are collected in the [go-dedup](https://github.com/go-dedup) collection as Go source, including **Bag of Words**, **Shingling** and **Hashing** (**MinHash** and **SimHash** to be exact). The **MinHash** and **SimHash** are the _latest_ algorithms, they both use a renovative technology called “[locality sensitive hashing (LSH)](https://en.wikipedia.org/wiki/Locality_sensitive_hashing)”.  They are used by search engines like mozilla and google currently, i.e. they are the _best_ open source algorithms _available_ at the present. `fsimilar` only uses two of them. Here is why.

### MinHash

[**MinHash**](http://en.wikipedia.org/wiki/MinHash) is an approach that is capable of using constant storage independent of the document length and producing a good estimate of our similarity measure. This approach distills down each document to a **fixed-size set of hashes** as a rough signature of this document. It was initially used in the [AltaVista](https://en.wikipedia.org/wiki/AltaVista) search engine to detect duplicate web pages and eliminate them from search results. It has also been applied in large-scale clustering problems, such as [clustering documents](https://en.wikipedia.org/wiki/Document_clustering) by the similarity of their sets of words.

### SimHash

While **MinHash** improves the complexity for comparing any two documents, we’re still faced with the problem of finding any (or all) duplicates of a query document given a set of known documents. Using this method, we have to do a linear traversal of all the known documents in O(n) time.

[**Simhashing**](https://en.wikipedia.org/wiki/SimHash) is a technique that helps us to overcome this problem. Given a set of input hashes, simhash produces a **single** hash with an interesting property — two similar sets of input hashes produce _similar_ resultant hashes, which distinctly different than most of hashing functions, which are chosen for the property that slightly different inputs yielding very-different outputs. This makes simhash unique and outstanding. [Google Crawler](https://en.wikipedia.org/wiki/Web_crawler) uses this algorithm to find near duplicate pages. Google uses the Minhash and LSH for [Google News](https://en.wikipedia.org/wiki/Google_News) personalization as well.

Of the above two [locality sensitive hashing](https://en.wikipedia.org/wiki/Locality_sensitive_hashing) techniques, `fsimilar` uses the latest and greates -- SimHash. But the [go-dedup](https://github.com/go-dedup) collection does include the MinHash algorithm as well.

The SimHash is one computation-level faster than MinHash, as it does not need do a _linear traversal_ of all the known documents -- it is doing _tree traversal_ instead. I.e., the algorithm is in O(log(n)) time instead of O(n) time. However, such super fast comes with a price, that it is no longer able to do the [documents clustering](https://en.wikipedia.org/wiki/Document_clustering) easily. This is why  `fsimilar` uses another algorithm to do the _precise_ [documents clustering](https://en.wikipedia.org/wiki/Document_clustering) -- the [Vector space model](https://en.wikipedia.org/wiki/Vector_space_model).

### Vector space model

The [Vector space model](https://en.wikipedia.org/wiki/Vector_space_model) is the widest accepted and used algorithm from the **Bag of Words** realm. It is also called “**term vector model**”, and it is

> an algebraic model for representing text documents (and any objects, in general) as vectors of identifiers, such as, for example, index terms. It is used in information filtering, information retrieval, indexing and relevancy rankings.

It is the foundation that the most famous **tf-idf** approach bases on. Vector space model has the following advantages:

- Simple model based on linear algebra
- Term weights not binary
- Allows computing a continuous degree of similarity between queries and documents
- Allows ranking documents according to their possible relevance
- Allows partial matching

It needs to do one-by-one comparison, so it is slower than SimHash in term of complexity. However, because it is doing one-by-one comparison, it is much accurate than SimHash in [_documents clustering_](https://en.wikipedia.org/wiki/Document_clustering), which is in fact what `fsimilar` is _actually doing_.

### Advantages of fsimilar over existing tools

There are many existing duplicate detection tools out there. Most of them do so by computing signature of the files, and then compare the signatures. The drawbacks of such approaches are obvious:

- Computing signature for files is slow. The _bigger_ the files are the _slower_ it will be.
- If to speed up the process, then extra disk space are necessary to store the signatures, which in turns may cause another problem that the signatures may be out of sync with the actual files.
- They only detect files that are exactly the same. I.e., same text file document backed up at different times, with slightly changes in between will _not_ be found as duplicates by such tools. Neither would be same files in different formats, `.pdf` vs. `.doc` vs `.epub` vs `.zip` vs `.rar`, or `.mp4`  vs. `.avi` vs `.mkv`, etc, etc.
- The exact duplicate of the same file are actually very rare in the file system. Similar files however are several magnitudes more.
- If the file names are different, such tool will be having a hard time dealing with it.

`fsimilar` on the other hand,

- only make guesses from the file names (and optionally file sizes). Thus it is several magnitudes faster than any of the signature creating/checking tools.
- it utilizes the state of the art technologies. So if used properly, it may be faster doing _similarity check_ than you do _file listing_ via `find /`, i.e., calling `find` to find all the files from the system root.
- because it is only guessing from such little information, it makes recommendation but not decisions.
- but it has many ways to make the decision-making easier, and the action to resolve the duplications easier too.

Moreover, there is one technique used by `fsimilar`that is not listed in above Moz article -- it can do **phonetic similarity checking** as well. This is why it is able to identify that “`Audio Book - The Grey Coloured Bunnie.mp3`” and “`ColoredGrayBunny.ogg`” are _maybe_ the same file! Because a [**phonetic algorithm**](https://en.wikipedia.org/wiki/Phonetic_algorithm) matches two different words with similar pronunciation to the same code, which allows phonetic similarity based word set comparison and indexing.

### Phonetic algorithms

- There are many [**phonetic algorithm**](https://en.wikipedia.org/wiki/Phonetic_algorithm) with various accuracy.
- The [go-dedup collection](https://github.com/go-dedup) includes all the more accurate ones.
- The `fsimilar` uses the [most accurate one](https://github.com/go-dedup/megophone) available in the roylty-free open-source world.

## Using fsimilar

### Online usage help

Online usage help that comes with `fsimilar`.

#### $ fsimilar
```sh
find/file similar
Version 0.9.0 built on 2017-09-10

Find similar files

Options:

  -h, --help            display help information
  -S, --size-given      size of the files in input as first field
  -Q, --query-size      query the file sizes from os
  -i, --input          *input from stdin or the given file (mandatory)
  -p, --phonetic        use phonetic as words for further error tolerant
  -F, --final           produce final output, the recommendations
  -c, --cp[=$FSIM_CP]   config path, path that hold all template files
  -v, --verbose         verbose mode (multiple -v increase the verbosity)

Commands:

  sim   Filter the input using simhash similarity check
  vec   Use Vector Space for similarity check
```

The full help at root level.

#### $ fsimilar sim
```sh
Filter the input using simhash similarity check

Usage:
  mlocate -i soccer | fsimilar sim -i

Options:

  -h, --help            display help information
  -S, --size-given      size of the files in input as first field
  -Q, --query-size      query the file sizes from os
  -i, --input          *input from stdin or the given file (mandatory)
  -p, --phonetic        use phonetic as words for further error tolerant
  -F, --final           produce final output, the recommendations
  -c, --cp[=$FSIM_CP]   config path, path that hold all template files
  -v, --verbose         verbose mode (multiple -v increase the verbosity)
  -d, --dist[=3]        the hamming distance of hashes within which to deem similar
```

This is sub-command `sim` level help, that filters the input using simhash similarity check.

#### $ fsimilar vec
```sh
Use Vector Space for similarity check

Usage:
  { mlocate -i soccer; mlocate -i football; } | fsimilar sim -i | fsimilar vec -i -S -Q -F

Options:

  -h, --help            display help information
  -S, --size-given      size of the files in input as first field
  -Q, --query-size      query the file sizes from os
  -i, --input          *input from stdin or the given file (mandatory)
  -p, --phonetic        use phonetic as words for further error tolerant
  -F, --final           produce final output, the recommendations
  -c, --cp[=$FSIM_CP]   config path, path that hold all template files
  -v, --verbose         verbose mode (multiple -v increase the verbosity)
  -t, --thr[=0.86]      the threshold above which to deem similar (0.8 = 80%)
```

This is sub-command `vec` level help, which uses Vector Space for similarity checking and reporting.

### Usage example

Let's take a look at how to use `fsimilar`. There are three test files in the `test/` directory, all generated with the same `find test/sim -type f` command. The only different between them is the actual file order in the list files. We'll use the first one for explanation in the following scenarios.

#### > test/sim.lstA
```sh
test/sim/Audio Book - The Grey Coloured Bunnie.mp3
test/sim/GNU - Python Standard Library (2001).rar
test/sim/PopupTest.java
test/sim/(eBook) GNU - Python Standard Library 2001.pdf
test/sim/Python Standard Library.zip
test/sim/GNU - 2001 - Python Standard Library.pdf
test/sim/LayoutTest.java
test/sim/ColoredGrayBunny.ogg
```

Again, this is one order that the files under the `test/sim/` directory may be fed into `fsimilar`.

#### $ fsimilar sim -i test/sim.lstA -d 12
```sh
### #1: test/sim/(eBook) GNU - Python Standard Library 2001.pdf

       1 test/sim/(eBook) GNU - Python Standard Library 2001.pdf
       1 test/sim/GNU - 2001 - Python Standard Library.pdf
       1 test/sim/GNU - Python Standard Library (2001).rar

### #2: test/sim/Python Standard Library.zip

       1 test/sim/Python Standard Library.zip
```

This shows how the result looks like using SimHash similar detection algorithm. We can see that all four Python books have been correctly identified as duplication candidates, but SimHash is just not able to group all four books into one cluster.

#### $ fsimilar vec -i test/sim.lstA -t 0.7
```sh
### #1: test/sim/GNU - Python Standard Library (2001).rar

       1 test/sim/GNU - Python Standard Library (2001).rar
       1 test/sim/(eBook) GNU - Python Standard Library 2001.pdf
       1 test/sim/Python Standard Library.zip
       1 test/sim/GNU - 2001 - Python Standard Library.pdf
```

This shows how the result looks like using the [Vector Space Model](https://en.wikipedia.org/wiki/Vector_space_model) similar detection algorithm. We can clearly see that all four books have been categorized into one _single_ cluster.

#### $ fsimilar vec -i test/sim.lstA -t 0.7 -p
```sh
### #1: test/sim/Audio Book - The Grey Coloured Bunnie.mp3

       1 test/sim/Audio Book - The Grey Coloured Bunnie.mp3
       1 test/sim/ColoredGrayBunny.ogg

### #2: test/sim/GNU - Python Standard Library (2001).rar

       1 test/sim/GNU - Python Standard Library (2001).rar
       1 test/sim/(eBook) GNU - Python Standard Library 2001.pdf
       1 test/sim/Python Standard Library.zip
       1 test/sim/GNU - 2001 - Python Standard Library.pdf
```

This is where `fsimilar` _shines_.

- This time it outputs another group, the "Audio Book" group.
- It understands sometimes the single word like "_ColoredGrayBunny_" needs to be split into individual words like "_Colored Gray Bunny_" before doing the similarity comparison.
- Although not a single word in the two books are the same, the `fsimilar` is able to identify that the two books  _maybe_ the same file, because the `-p` tells it to do the [**phonetic similarity checking**](https://en.wikipedia.org/wiki/Phonetic_algorithm).

#### alternative invocation

Instead of specifying the input after `-i` (e.g., `-i test/sim.lstA`), `fsimilar` can also accept input from standard input as well:

    find test/sim -type f | fsimilar vec -i -t ...

It'll be the same as `fsimilar vec -i test/sim.lstA -t ...`.

### Using `fsimilar`

Now, put all above into practice, to use `fsimilar` with a concrete example:

    fsimilar vec -i test/test1.lst -S -F

The `-F`, i.e. `--final`, tells `fsimilar` to _produce final output_.

#### > head -58 test/shell_test.ref
```sh
### #1

100  5443429 Soccer Positions Explained.mp4	./Soccer Tips & Soccer Advice/
 87  5443429 22 - Soccer Positions Explained.mp4	./How to improve Soccer conditioning & Soccer fitness/

### #2

100  5973267 Football Dribbling Tips - How to dribble.mp4	./Soccer Tips & Soccer Advice/
100  5973267 Football Dribbling Tips - How to dribble.mp4	./Top Soccer Training Videos/

### #3

100  4455628 Soccer Training Guide.webm	./Soccer Tips & Soccer Advice/
100  4455628 Soccer Training Guide.webm	./How to improve Soccer passing & receiving skills/
100  4455628 Soccer Training Guide.webm	./How to improve Soccer ball control skills/
100  4455628 Soccer Training Guide.webm	./Youth Soccer Training Drills/
100  4455628 Soccer Training Guide.webm	./How to improve Soccer shooting skills & finishing/
100  4455628 Soccer Training Guide.webm	./Top Soccer Training Videos/
 87  4455628 32 - Soccer Training Guide.webm	./How to improve Soccer conditioning & Soccer fitness/
 87  4455628 20 - Soccer Training Guide.webm	./At Home Soccer Training Drills/
100  4455628 Soccer Training Guide.webm	./How to improve Soccer dribbling skills and fast footwork/

### #4

100  5805224 Soccer Tips - How to play midfielder.mp4	./Soccer Tips & Soccer Advice/
100  5805224 Soccer Tips - How to play midfielder.mp4	./Top Soccer Training Videos/

### #5

100  4090343 Football Tips - How to play football.mp4	./Soccer Tips & Soccer Advice/
100  4090343 Football Tips - How to play football.mp4	./Top Soccer Training Videos/

### #6

100  5131092 Soccer Shooting Tips - How to Shoot a Soccer Ball.mp4	./Soccer Tips & Soccer Advice/
100  5131092 Soccer Shooting Tips - How to Shoot a Soccer Ball.mp4	./Top Soccer Training Videos/

### #7

100 13599728 Football Passing Drills.mp4	./How to improve Soccer passing & receiving skills/
100 13599728 Football Passing Drills.mp4	./Youth Soccer Training Drills/

### #8

100 13295374 Soccer Passing Drills For Youth.mp4	./How to improve Soccer passing & receiving skills/
100 13295374 Soccer Passing Drills For Youth.mp4	./Youth Soccer Training Drills/

### #9

100 11827630 Soccer Passing Drills.mp4	./How to improve Soccer passing & receiving skills/
100 11827630 Soccer Passing Drills.mp4	./Youth Soccer Training Drills/

### #10

100  7997569 Soccer Ball Control Drills.mp4	./How to improve Soccer passing & receiving skills/
100  7997569 Soccer Ball Control Drills.mp4	./Youth Soccer Training Drills/
100  7997569 Soccer Ball Control Drills.mp4	./Top Soccer Training Videos/
```

The above lists first 10 group of similar files, produced via the above _produce final output_ command.

- The **first** column is the _confident level_, with 100 meaning 100% confident.
- The **second** column is the _size_ of the file, which is an important factor in determining if the two files are the same.
- The **third** column is the _name_ of the file, from which `fsimilar` tried to make the guess for the similarity.
- The **forth** column is the _directory name_ where the file resides. If the file name and file size are the same, they must be in different directory.
- The *first listing* is comparing to itself, thus its confident level will *always* be 100%.

In the first group, we can see that the confident level for the second file is only 87%, because it's name is different than the first one. Group 3 has files with confident level at 87% as well. How to use this confident level will be explained next.

### The recommendations and decisions

When instructed to _produce final output_, `fsimilar` will also generate a bunch of shell commands under the TEMP directory. See [test/shell_test.sh](test/shell_test.sh) for details. Here are the shell commands for `ln` and `mv`, i.e.,

- overwriting the similar files by a _symlink_ to the original, or
- deleting the similar files by moving them into a "_trash_" folder.
- for _straight deleting_ without saving, take a look at the [`rm` shell command](test/shell_rm.tmpl.sh).

#### > head -18 test/shell_ln.tmpl.sh
```sh
### #1: ./Soccer Tips & Soccer Advice/Soccer Positions Explained.mp4

  [ "$FSIM_MIN" ] && [  87 -ge "$FSIM_MIN" ] && $FSIM_SHOW ln -sf '../Soccer Tips & Soccer Advice/Soccer Positions Explained.mp4' './How to improve Soccer conditioning & Soccer fitness/22 - Soccer Positions Explained.mp4'

### #2: ./Soccer Tips & Soccer Advice/Football Dribbling Tips - How to dribble.mp4

  [ "$FSIM_MIN" ] && [ 100 -ge "$FSIM_MIN" ] && $FSIM_SHOW ln -sf '../Soccer Tips & Soccer Advice/Football Dribbling Tips - How to dribble.mp4' './Top Soccer Training Videos/Football Dribbling Tips - How to dribble.mp4'

### #3: ./Soccer Tips & Soccer Advice/Soccer Training Guide.webm

  [ "$FSIM_MIN" ] && [ 100 -ge "$FSIM_MIN" ] && $FSIM_SHOW ln -sf '../Soccer Tips & Soccer Advice/Soccer Training Guide.webm' './How to improve Soccer passing & receiving skills/Soccer Training Guide.webm'
  [ "$FSIM_MIN" ] && [ 100 -ge "$FSIM_MIN" ] && $FSIM_SHOW ln -sf '../Soccer Tips & Soccer Advice/Soccer Training Guide.webm' './How to improve Soccer ball control skills/Soccer Training Guide.webm'
  [ "$FSIM_MIN" ] && [ 100 -ge "$FSIM_MIN" ] && $FSIM_SHOW ln -sf '../Soccer Tips & Soccer Advice/Soccer Training Guide.webm' './Youth Soccer Training Drills/Soccer Training Guide.webm'
  [ "$FSIM_MIN" ] && [ 100 -ge "$FSIM_MIN" ] && $FSIM_SHOW ln -sf '../Soccer Tips & Soccer Advice/Soccer Training Guide.webm' './How to improve Soccer shooting skills & finishing/Soccer Training Guide.webm'
  [ "$FSIM_MIN" ] && [ 100 -ge "$FSIM_MIN" ] && $FSIM_SHOW ln -sf '../Soccer Tips & Soccer Advice/Soccer Training Guide.webm' './Top Soccer Training Videos/Soccer Training Guide.webm'
  [ "$FSIM_MIN" ] && [  87 -ge "$FSIM_MIN" ] && $FSIM_SHOW ln -sf '../Soccer Tips & Soccer Advice/Soccer Training Guide.webm' './How to improve Soccer conditioning & Soccer fitness/32 - Soccer Training Guide.webm'
  [ "$FSIM_MIN" ] && [  87 -ge "$FSIM_MIN" ] && $FSIM_SHOW ln -sf '../Soccer Tips & Soccer Advice/Soccer Training Guide.webm' './At Home Soccer Training Drills/20 - Soccer Training Guide.webm'
  [ "$FSIM_MIN" ] && [ 100 -ge "$FSIM_MIN" ] && $FSIM_SHOW ln -sf '../Soccer Tips & Soccer Advice/Soccer Training Guide.webm' './How to improve Soccer dribbling skills and fast footwork/Soccer Training Guide.webm'
```

#### > head -18 test/shell_mv.tmpl.sh
```sh
### #1: ./Soccer Tips & Soccer Advice/Soccer Positions Explained.mp4

  # [1] Extension '.mp4' of './How to improve Soccer conditioning & Soccer fitness/22 - Soccer Positions Explained.mp4' same as original, retained

### #2: ./Soccer Tips & Soccer Advice/Football Dribbling Tips - How to dribble.mp4

  # [1] Extension '.mp4' of './Top Soccer Training Videos/Football Dribbling Tips - How to dribble.mp4' same as original, retained

### #3: ./Soccer Tips & Soccer Advice/Soccer Training Guide.webm

  # [1] Extension '.webm' of './How to improve Soccer passing & receiving skills/Soccer Training Guide.webm' same as original, retained
  # [2] Extension '.webm' of './How to improve Soccer ball control skills/Soccer Training Guide.webm' same as original, retained
  # [3] Extension '.webm' of './Youth Soccer Training Drills/Soccer Training Guide.webm' same as original, retained
  # [4] Extension '.webm' of './How to improve Soccer shooting skills & finishing/Soccer Training Guide.webm' same as original, retained
  # [5] Extension '.webm' of './Top Soccer Training Videos/Soccer Training Guide.webm' same as original, retained
  # [6] Extension '.webm' of './How to improve Soccer conditioning & Soccer fitness/32 - Soccer Training Guide.webm' same as original, retained
  # [7] Extension '.webm' of './At Home Soccer Training Drills/20 - Soccer Training Guide.webm' same as original, retained
  # [8] Extension '.webm' of './How to improve Soccer dribbling skills and fast footwork/Soccer Training Guide.webm' same as original, retained
```

#### The approach

The above two listings list what/how `fsimilar` recommends to do, for the first three group of similar files. The `mv` shell command file is all-comments (1). If you try with the `test/sim.lstA` instead, the command/comment will actually be reversed for the `ln` and `mv` shell commands. In reality, maybe both will contain some commands and some comments.

(1) The reason `mv` shell command file is all-comments is for travis CI testing to pass. If using the `shell_mv.tmpl.real` instead as the template, the output will be different in each run as the "_trash_" folder is different in each run (so as not to loose any files by accident).

As stressed before, `fsimilar` only makes recommendation, not decisions. Now it is time for you, the end user, to determine how to make the decisions. The general rule of thumb is,

- Look carefully at the [similarity report](test/shell_test.ref) to see if there is anything that are recommended wrongly.
- Increase the reporting threshold on the `fsimilar` command line so that the  wrongly recommended are minimum.
- Find according to the `### #XX:` from the [similarity report](test/shell_test.ref), into the shell command files and remove those wrong commands.
- Use the `mv` shell command file to make a backup of those similar files first.
- Then use the `ln` shell command file to finally re-create the similar files by a symlink to the original.
- Set `FSIM_SHOW` to `echo` to dry run the shell command before actually doing it (when `FSIM_SHOW=''`)
- Set `FSIM_MIN` to `100` to deal with sure-duplicate files first, then lower it value gradually to deal with them in "waves". Or remove the 100% duplicates first so as to deal with a much shorter list later manually.
- Note that both `FSIM_SHOW` and `FSIM_MIN` need to be `export`ed for the shell script to pick them up.
- When the reporting threshold are set to too low to catch a certain file, manually copy & paste that specific command into console instead of dealing with the shell script as a whole.

## Binary releases

_Coming soon_.

# Download/Install

## Download binaries

- The latest binary executables are available under  
https://bintray.com/suntong/bin/fsimilar#files/fsimilar  
as the result of the Continuous-Integration process.
- I.e., they are built right from the source code during every git commit automatically by [travis-ci](https://travis-ci.org/).
- Pick & choose the binary executable that suits your OS and its architecture. E.g., for Linux, it would most probably be the `fsimilar_linux_VER_amd64` file. If your OS and its architecture is not available in the download list, please let me know and I'll add it.
- You may want to rename it to a shorter name instead, e.g., `fsimilar`, after downloading it.


## Debian package

Available at https://dl.bintray.com/suntong/deb.

```
echo "deb [trusted=yes] https://dl.bintray.com/suntong/deb all main" | sudo tee /etc/apt/sources.list.d/suntong-debs.list
sudo apt-get update

sudo chmod 644 /etc/apt/sources.list.d/suntong-debs.list
apt-cache policy fsimilar

sudo apt-get install -y fsimilar
```


## Install Source

To install the source code instead:

```
go get github.com/go-dedup/fsimilar
```

Note that, this project is also a showcase of the [Go wire-framing solution](https://github.com/go-easygen/wireframe), that wire-frames Go cli based projects from start to finish, including but not limited to,

- [Create Repository in Github](https://github.com/go-easygen/wireframe#github-repo-create---create-repository-in-github)
- [Command line flag handling code auto-generation](https://github.com/go-easygen/wireframe#command-line-flag-handling-code-auto-generation)
- [Creating binary releases on bintray](https://github.com/go-easygen/wireframe#binary-releases)

# Similar Projects

This project is based on the [File-FindSimilars](https://metacpan.org/release/File-FindSimilars) Perl module, which is released on [CPAN](https://metacpan.org/) more than a decade ago.

# Credits

See the individual Go source modules that this project is based on.

# Author(s) & Contributor(s)

Tong SUN  
![suntong from cpan.org](https://img.shields.io/badge/suntong-%40cpan.org-lightgrey.svg "suntong from cpan.org")
