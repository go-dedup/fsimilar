
# {{.Name}}

{{render "license/shields" . "License" "MIT"}}
{{template "badge/godoc" .}}
{{template "badge/goreport" .}}
{{template "badge/travis" .}}

## {{toc 5}}

## {{.Name}} - find/file similar

`{{.Name}}` find similar file very fast -- it can spot similar files within the file system very quickly.

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
- The `{{.Name}}` uses the [most accurate one](https://github.com/go-dedup/megophone) available in the roylty-free open-source world.

## Using {{.Name}}

### Online usage help

Online usage help that comes with `{{.Name}}`.

#### $ {{exec "fsimilar" | color "sh"}}

The full help at root level.

#### $ {{shell "fsimilar sim" | color "sh"}}

This is sub-command `sim` level help, that filters the input using simhash similarity check.

#### $ {{shell "fsimilar vec" | color "sh"}}

This is sub-command `vec` level help, which uses Vector Space for similarity checking and reporting.

### Usage example

Let's take a look at how to use `{{.Name}}`. There are three test files in the `test/` directory, all generated with the same `find test/sim -type f` command. The only different between them is the actual file order in the list files. We'll use the first one for explanation in the following scenarios.

#### > {{cat "test/sim.lstA" | color "sh"}}

Again, this is one order that the files under the `test/sim/` directory may be fed into `{{.Name}}`.

#### $ {{shell "fsimilar sim -i test/sim.lstA -d 12" | color "sh"}}

This shows how the result looks like using SimHash similar detection algorithm. We can see that all four Python books have been correctly identified as duplication candidates, but SimHash is just not able to group all four books into one cluster.

#### $ {{shell "fsimilar vec -i test/sim.lstA -t 0.7" | color "sh"}}

This shows how the result looks like using the [Vector Space Model](https://en.wikipedia.org/wiki/Vector_space_model) similar detection algorithm. We can clearly see that all four books have been categorized into one _single_ cluster.

#### $ {{shell "fsimilar vec -i test/sim.lstA -t 0.7 -p" | color "sh"}}

This is where `{{.Name}}` _shines_.

- This time it outputs another group, the "Audio Book" group.
- It understands sometimes the single word like "_ColoredGrayBunny_" needs to be split into individual words like "_Colored Gray Bunny_" before doing the similarity comparison.
- Although not a single word in the two books are the same, the `{{.Name}}` is able to identify that the two books  _maybe_ the same file, because the `-p` tells it to do the [**phonetic similarity checking**](https://en.wikipedia.org/wiki/Phonetic_algorithm).

#### alternative invocation

Instead of specifying the input after `-i` (e.g., `-i test/sim.lstA`), `{{.Name}}` can also accept input from standard input as well:

    find test/sim -type f | fsimilar vec -i -t ...

It'll be the same as `fsimilar vec -i test/sim.lstA -t ...`.

### Using `{{.Name}}`

Now, put all above into practice, to use `{{.Name}}` with a concrete example:

    fsimilar vec -i test/test1.lst -S -F

The `-F`, i.e. `--final`, tells `{{.Name}}` to _produce final output_.

#### > {{shell "head -58 test/shell_test.ref" | color "sh"}}

The above lists first 10 group of similar files, produced via the above _produce final output_ command.

- The **first** column is the _confident level_, with 100 meaning 100% confident.
- The **second** column is the _size_ of the file, which is an important factor in determining if the two files are the same.
- The **third** column is the _name_ of the file, from which `{{.Name}}` tried to make the guess for the similarity.
- The **forth** column is the _directory name_ where the file resides. If the file name and file size are the same, they must be in different directory.
- The *first listing* is comparing to itself, thus its confident level will *always* be 100%.

In the first group, we can see that the confident level for the second file is only 87%, because it's name is different than the first one. Group 3 has files with confident level at 87% as well. How to use this confident level will be explained next.

### The recommendations and decisions

When instructed to _produce final output_, `{{.Name}}` will also generate a bunch of shell commands under the TEMP directory. See [test/shell_test.sh](test/shell_test.sh) for details. Here are the shell commands for `ln` and `mv`, i.e.,

- overwriting the similar files by a _symlink_ to the original, or
- deleting the similar files by moving them into a "_trash_" folder.
- for _straight deleting_ without saving, take a look at the [`rm` shell command](test/shell_rm.tmpl.sh).

#### > {{shell "head -18 test/shell_ln.tmpl.sh" | color "sh"}}

#### > {{shell "head -18 test/shell_mv.tmpl.sh" | color "sh"}}

#### The approach

The above two listings list what/how `{{.Name}}` recommends to do, for the first three group of similar files. The `mv` shell command file is all-comments (1). If you try with the `test/sim.lstA` instead, the command/comment will actually be reversed for the `ln` and `mv` shell commands. In reality, maybe both will contain some commands and some comments.

(1) The reason `mv` shell command file is all-comments is for travis CI testing to pass. If using the `shell_mv.tmpl.real` instead as the template, the output will be different in each run as the "_trash_" folder is different in each run (so as not to loose any files by accident).

As stressed before, `{{.Name}}` only makes recommendation, not decisions. Now it is time for you, the end user, to determine how to make the decisions. The general rule of thumb is,

- Look carefully at the [similarity report](test/shell_test.ref) to see if there is anything that are recommended wrongly.
- Increase the reporting threshold on the `{{.Name}}` command line so that the  wrongly recommended are minimum.
- Find according to the `### #XX:` from the [similarity report](test/shell_test.ref), into the shell command files and remove those wrong commands.
- Use the `mv` shell command file to make a backup of those similar files first.
- Then use the `ln` shell command file to finally re-create the similar files by a symlink to the original.
- Set `FSIM_SHOW` to `echo` to dry run the shell command before actually doing it (when `FSIM_SHOW=''`)
- Set `FSIM_MIN` to `100` to deal with sure-duplicate files first, then lower it value gradually to deal with them in "waves". Or remove the 100% duplicates first so as to deal with a much shorter list later manually.
- Note that both `FSIM_SHOW` and `FSIM_MIN` need to be `export`ed for the shell script to pick them up.
- When the reporting threshold are set to too low to catch a certain file, manually copy & paste that specific command into console instead of dealing with the shell script as a whole.

# Download/Install

## Download binaries

- The latest binary executables are available under  
https://bintray.com/suntong/bin/{{.Name}}/latest  
as the result of the Continuous-Integration process.
- I.e., they are built right from the source code during _every_ git commit _automatically_ by [travis-ci](https://travis-ci.org/).
- Pick & choose the binary executable that suits your OS and its architecture. E.g., for Linux, it would most probably be the `{{.Name}}-linux-amd64` file. If your OS and its architecture is not available in the download list, please let me know and I'll add it.
- You may want to rename it to a shorter name instead, e.g., `{{.Name}}`, after downloading it.


## Debian package

Debian package _repo_ is available at https://dl.bintray.com/suntong/deb.
The _browse-able_ repo view is at https://bintray.com/suntong/deb.

```
echo "deb [trusted=yes] https://dl.bintray.com/suntong/deb all main" | sudo tee /etc/apt/sources.list.d/suntong-debs.list
sudo apt-get update

sudo chmod 644 /etc/apt/sources.list.d/suntong-debs.list
apt-cache policy {{.Name}}

sudo apt-get install -y {{.Name}}
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

