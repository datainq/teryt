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
go test -test.v -test.bench ^BenchmarkSearch_Search.*$ -test.run ^$
goos: linux
goarch: amd64
pkg: github.com/datainq/teryt/search

BenchmarkSearch_SearchOld
BenchmarkSearch_SearchOld-8   	      37	  33642952 ns/op
BenchmarkSearch_Search
BenchmarkSearch_Search-8      	     152	   7784588 ns/op
BenchmarkSearch_SearchV3
BenchmarkSearch_SearchV3/parallel-1
BenchmarkSearch_SearchV3/parallel-1-8         	     150	   7900379 ns/op
BenchmarkSearch_SearchV3/parallel-2
BenchmarkSearch_SearchV3/parallel-2-8         	     277	   4286875 ns/op
BenchmarkSearch_SearchV3/parallel-3
BenchmarkSearch_SearchV3/parallel-3-8         	     321	   3643391 ns/op
BenchmarkSearch_SearchV3/parallel-4
BenchmarkSearch_SearchV3/parallel-4-8         	     330	   3492188 ns/op
BenchmarkSearch_SearchV3/parallel-5
BenchmarkSearch_SearchV3/parallel-5-8         	     382	   3085371 ns/op
BenchmarkSearch_SearchV3/parallel-6
BenchmarkSearch_SearchV3/parallel-6-8         	     429	   2752282 ns/op
BenchmarkSearch_SearchV3/parallel-7
BenchmarkSearch_SearchV3/parallel-7-8         	     493	   2352805 ns/op
BenchmarkSearch_SearchV3/parallel-8
BenchmarkSearch_SearchV3/parallel-8-8         	     490	   2433284 ns/op
PASS
ok  	github.com/datainq/teryt/search	27.525s
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
BenchmarkAll/ASCII/agniva-8  	 3711177	       322 ns/op
BenchmarkAll/ASCII/m1ome
BenchmarkAll/ASCII/m1ome-8   	 2355360	       510 ns/op
BenchmarkAll/ASCII/arbovm
BenchmarkAll/ASCII/arbovm-8  	 2955633	       404 ns/op
BenchmarkAll/ASCII/dgryski
BenchmarkAll/ASCII/dgryski-8 	 2792784	       416 ns/op
BenchmarkAll/ASCII/datainq
BenchmarkAll/ASCII/datainq-8 	 3908599	       305 ns/op

BenchmarkAll/French
BenchmarkAll/French/agniva
BenchmarkAll/French/agniva-8 	 2569628	       465 ns/op
BenchmarkAll/French/m1ome
BenchmarkAll/French/m1ome-8  	 3487110	       343 ns/op
BenchmarkAll/French/arbovm
BenchmarkAll/French/arbovm-8 	 2440542	       487 ns/op
BenchmarkAll/French/dgryski
BenchmarkAll/French/dgryski-8         	 2262538	       522 ns/op
BenchmarkAll/French/datainq
BenchmarkAll/French/datainq-8         	 2607481	       445 ns/op

BenchmarkAll/Nordic
BenchmarkAll/Nordic/agniva
BenchmarkAll/Nordic/agniva-8          	 1276735	       929 ns/op
BenchmarkAll/Nordic/m1ome
BenchmarkAll/Nordic/m1ome-8           	 2263801	       534 ns/op
BenchmarkAll/Nordic/arbovm
BenchmarkAll/Nordic/arbovm-8          	 1281937	       947 ns/op
BenchmarkAll/Nordic/dgryski
BenchmarkAll/Nordic/dgryski-8         	 1000000	      1022 ns/op
BenchmarkAll/Nordic/datainq
BenchmarkAll/Nordic/datainq-8         	 1399813	       866 ns/op

BenchmarkAll/Tibetan
BenchmarkAll/Tibetan/agniva
BenchmarkAll/Tibetan/agniva-8         	 1438171	       833 ns/op
BenchmarkAll/Tibetan/m1ome
BenchmarkAll/Tibetan/m1ome-8          	 1711886	       702 ns/op
BenchmarkAll/Tibetan/arbovm
BenchmarkAll/Tibetan/arbovm-8         	 1373340	       890 ns/op
BenchmarkAll/Tibetan/dgryski
BenchmarkAll/Tibetan/dgryski-8        	 1302441	       905 ns/op
BenchmarkAll/Tibetan/datainq
BenchmarkAll/Tibetan/datainq-8        	 1473806	       802 ns/op

BenchmarkAll/Polish
BenchmarkAll/Polish/agniva
BenchmarkAll/Polish/agniva-8          	 3482793	       341 ns/op
BenchmarkAll/Polish/m1ome
BenchmarkAll/Polish/m1ome-8           	 2719704	       437 ns/op
BenchmarkAll/Polish/arbovm
BenchmarkAll/Polish/arbovm-8          	 3161550	       379 ns/op
BenchmarkAll/Polish/dgryski
BenchmarkAll/Polish/dgryski-8         	 3030890	       397 ns/op
BenchmarkAll/Polish/datainq
BenchmarkAll/Polish/datainq-8         	 3609500	       334 ns/op
PASS

Process finished with exit code 0
```