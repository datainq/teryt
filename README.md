# teryt ![CI](https://github.com/datainq/teryt/workflows/CI/badge.svg)

[TERYT](https://pl.wikipedia.org/wiki/TERYT) - Krajowy Rejestr Urzędowy Podziału Terytorialnego Kraju
is an Polish official registry of country administrative division. 

Source: http://eteryt.stat.gov.pl/eTeryt/rejestr_teryt/udostepnianie_danych/baza_teryt/uzytkownicy_indywidualni/pobieranie/pliki_pelne_struktury.aspx

The registry has four parts:
 * TERC - identifiers and names of administrative units,
 * SIMC - identifiers nad names of places,
 * BREC - statistical regions and census areas,
 * NOBC - (it includes ULIC), address identification of streets, estate, buildings and apartments.

This library provides parser for the TERC, SIMC, ULIC
datasets.

The full data files are available at official GUS website:
http://eteryt.stat.gov.pl/eTeryt/rejestr_teryt/udostepnianie_danych/baza_teryt/uzytkownicy_indywidualni/pobieranie/pobieranie.aspx?contrast=default

## Install

```
go get github.com/datainq/teryt
```

## Usage

You need to download data files. After loading you can work on it.

Check example search program in `cmd/search.go`. 
You can run it by executing in a main repo directory:

```
> go run cmd/search.go -simc data/SIMC_Urzedowy_2020-09-26.zip -terc data/TERC_Urzedowy_2020-09-26.zip

INFO[0000] terc: data/TERC_Urzedowy_2020-09-26.zip
INFO[0000] simc: data/SIMC_Urzedowy_2020-09-26.zip
>Olsztyn
search for "olsztyn"
Results: (24.767138ms)
0. dist 0: WARMIŃSKO-MAZURSKIE Olsztyn (miasto na prawach powiatu)
1. dist 0: ŚLĄSKIE częstochowski Olsztyn (gmina wiejska)
2. dist 0: ŚLĄSKIE częstochowski Olsztyn Olsztyn (wieś)
3. dist 1: WIELKOPOLSKIE wolsztyński Wolsztyn (gmina miejsko-wiejska)
4. dist 1: WIELKOPOLSKIE wolsztyński Wolsztyn (miasto)
5. dist 1: WIELKOPOLSKIE wolsztyński Wolsztyn (obszar wiejski)
6. dist 1: LUBELSKIE bialski Rokitno Olszyn (wieś)
7. dist 1: WIELKOPOLSKIE wolsztyński Wolsztyn Wolsztyn (miasto)
8. dist 1: WIELKOPOLSKIE czarnkowsko-trzcianecki Wieleń Folsztyn (wieś)
9. dist 2: DOLNOŚLĄSKIE lubański Olszyna (gmina miejsko-wiejska)
>
```

## Benchmarks

```
go test -test.v -test.bench ^BenchmarkSearch_Search.*$ -test.run ^$ -benchtime=5s
goos: linux
goarch: amd64
pkg: github.com/datainq/teryt/search
BenchmarkSearch_SearchOld
BenchmarkSearch_SearchOld-8   	     172	  34626860 ns/op	 1804163 B/op	  114175 allocs/op
BenchmarkSearch_Search
BenchmarkSearch_Search-8      	     770	   7788946 ns/op	  923129 B/op	   57257 allocs/op
BenchmarkSearch_SearchV3
BenchmarkSearch_SearchV3/parallel-1
BenchmarkSearch_SearchV3/parallel-1-8         	     901	   6752492 ns/op	     690 B/op	      22 allocs/op
BenchmarkSearch_SearchV3/parallel-2
BenchmarkSearch_SearchV3/parallel-2-8         	    1743	   3525559 ns/op	    1055 B/op	      28 allocs/op
BenchmarkSearch_SearchV3/parallel-3
BenchmarkSearch_SearchV3/parallel-3-8         	    2450	   2430533 ns/op	    1412 B/op	      34 allocs/op
BenchmarkSearch_SearchV3/parallel-4
BenchmarkSearch_SearchV3/parallel-4-8         	    3187	   2001791 ns/op	    1778 B/op	      40 allocs/op
BenchmarkSearch_SearchV3/parallel-5
BenchmarkSearch_SearchV3/parallel-5-8         	    2359	   2427669 ns/op	    2154 B/op	      46 allocs/op
BenchmarkSearch_SearchV3/parallel-6
BenchmarkSearch_SearchV3/parallel-6-8         	    2626	   2103522 ns/op	    2491 B/op	      52 allocs/op
BenchmarkSearch_SearchV3/parallel-7
BenchmarkSearch_SearchV3/parallel-7-8         	    2953	   1835657 ns/op	    2869 B/op	      58 allocs/op
BenchmarkSearch_SearchV3/parallel-8
BenchmarkSearch_SearchV3/parallel-8-8         	    3070	   1693271 ns/op	    3211 B/op	      64 allocs/op
PASS
ok  	github.com/datainq/teryt/search	84.667s
```

## Benchmark of levenstein function

```
/tmp/___BenchmarkAll_in_github_com_datainq_teryt_search -test.v -test.bench ^\QBenchmarkAll\E$ -test.run ^$
goos: linux
goarch: amd64
pkg: github.com/datainq/teryt/search

BenchmarkAll

BenchmarkAll/ASCII
BenchmarkAll/ASCII/agniva
BenchmarkAll/ASCII/agniva-8  	 3696501	       319 ns/op
BenchmarkAll/ASCII/m1ome
BenchmarkAll/ASCII/m1ome-8   	 2402422	       501 ns/op
BenchmarkAll/ASCII/arbovm
BenchmarkAll/ASCII/arbovm-8  	 2978025	       401 ns/op
BenchmarkAll/ASCII/dgryski
BenchmarkAll/ASCII/dgryski-8 	 2841715	       419 ns/op
BenchmarkAll/ASCII/datainq
BenchmarkAll/ASCII/datainq-8 	 3842545	       306 ns/op
BenchmarkAll/ASCII/datainq2
BenchmarkAll/ASCII/datainq2-8         	 4183891	       283 ns/op

BenchmarkAll/Polish
BenchmarkAll/Polish/agniva
BenchmarkAll/Polish/agniva-8          	 3356847	       357 ns/op
BenchmarkAll/Polish/m1ome
BenchmarkAll/Polish/m1ome-8           	 2727634	       436 ns/op
BenchmarkAll/Polish/arbovm
BenchmarkAll/Polish/arbovm-8          	 3106969	       378 ns/op
BenchmarkAll/Polish/dgryski
BenchmarkAll/Polish/dgryski-8         	 3026116	       396 ns/op
BenchmarkAll/Polish/datainq
BenchmarkAll/Polish/datainq-8         	 3573074	       336 ns/op
BenchmarkAll/Polish/datainq2
BenchmarkAll/Polish/datainq2-8        	 3883668	       304 ns/op

BenchmarkAll/French
BenchmarkAll/French/agniva
BenchmarkAll/French/agniva-8          	 2499796	       481 ns/op
BenchmarkAll/French/m1ome
BenchmarkAll/French/m1ome-8           	 3454462	       340 ns/op
BenchmarkAll/French/arbovm
BenchmarkAll/French/arbovm-8          	 2388702	       489 ns/op
BenchmarkAll/French/dgryski
BenchmarkAll/French/dgryski-8         	 2227182	       544 ns/op
BenchmarkAll/French/datainq
BenchmarkAll/French/datainq-8         	 2646777	       462 ns/op
BenchmarkAll/French/datainq2
BenchmarkAll/French/datainq2-8        	 2740575	       428 ns/op

BenchmarkAll/Nordic
BenchmarkAll/Nordic/agniva
BenchmarkAll/Nordic/agniva-8          	 1236994	       986 ns/op
BenchmarkAll/Nordic/m1ome
BenchmarkAll/Nordic/m1ome-8           	 2088196	       520 ns/op
BenchmarkAll/Nordic/arbovm
BenchmarkAll/Nordic/arbovm-8          	 1243872	       952 ns/op
BenchmarkAll/Nordic/dgryski
BenchmarkAll/Nordic/dgryski-8         	 1000000	      1044 ns/op
BenchmarkAll/Nordic/datainq
BenchmarkAll/Nordic/datainq-8         	 1356010	       877 ns/op
BenchmarkAll/Nordic/datainq2
BenchmarkAll/Nordic/datainq2-8        	 1271364	       881 ns/op

BenchmarkAll/Tibetan
BenchmarkAll/Tibetan/agniva
BenchmarkAll/Tibetan/agniva-8         	 1413631	       873 ns/op
BenchmarkAll/Tibetan/m1ome
BenchmarkAll/Tibetan/m1ome-8          	 1459371	       750 ns/op
BenchmarkAll/Tibetan/arbovm
BenchmarkAll/Tibetan/arbovm-8         	 1330542	       961 ns/op
BenchmarkAll/Tibetan/dgryski
BenchmarkAll/Tibetan/dgryski-8        	 1212814	      1022 ns/op
BenchmarkAll/Tibetan/datainq
BenchmarkAll/Tibetan/datainq-8        	 1391295	       867 ns/op
BenchmarkAll/Tibetan/datainq2
BenchmarkAll/Tibetan/datainq2-8       	 1497633	       811 ns/op
PASS

Process finished with exit code 0
```