# WP-SEA

Projekt na zaliczenie Akademii Programowania WP PJATK
## Rozpoczęcie
### Instalacja
Aby pobrać ten program należy wykonać komendę:
```bash
  git clone https://github.com/joohnes/wp-sea
```
### Start

Aby uruchomić klienta, należy wykonać podaną niżej komendę w głównym folderze:

```bash
  go run main.go
```



## Jak zagrać w statki?
Po włączeniu aplikacji, klient da nam komunikat o nawiązywaniu połączenia z serwerem.
Po pomyślnym połączeniu pokaże nam się plansza razem z opisem naszej postaci i postaci przeciwnika

Następnie pokaże nam się komunikat o wpisaniu koordynatów naszego ataku.
Koordynaty powinny składać się z dwóch lub trzech znaków, pierwszy to litera od A do J, a drugi (czasami i trzeci) to liczba od 1 do 10.
W przypadku błędnego koordynatu program poinformuje o błędzie i poprosi jeszcze raz o wprowadzenie.

Przykładowe koordynaty:
```terminal
  A4, G10, J7, D10, F1
```

Jeśli chcemy się poddać, zamiast koordynatów możemy wpisać `quit`, co spowoduje wysłanie do serwera komunikatu o poddaniu i zakończeniu działania programu.
