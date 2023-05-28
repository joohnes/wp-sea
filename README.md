# WP-SEA

Projekt na zaliczenie Akademii Programowania WP PJATK
## Rozpoczęcie
### Start

Aby uruchomić klienta, należy wykonać podaną niżej komendę w głównym folderze:

```bash
  go run main.go
```



## Jak zagrać w statki?
Po włączeniu aplikacji, klient poprosi nas o podanie naszego pseudonimu i opisu. Następnie wyświetli się menu.\
Opcję wybieramy wpisując odpowiedni numer i wciskając klawisz 'enter'.

### Walka
Po wybraniu opcji z menu, klient połączy się z serwerem i przedstawi nam planszę.
W lewym górnym rogu mamy komunikat jak opuścić rozgrywkę (jest to również poddanie partii), i ile czasu nam zostało.
W prawym górnym rogu ukazują nam się nasze statystyki, ile razy strzeliliśmy, ile trafiliśmy, oraz procentowa celność.

Koordynaty wybieramy poprzez kliknięcie w wybrane miejsce lewym przyciskiem myszy.
- Klient nie pozwoli nam strzelić podczas trwania tury przeciwnika
- Klient nie pozwoli nam strzelić w miejsce, do którego już strzelaliśmy

