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