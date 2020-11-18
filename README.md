![Go](https://github.com/MasterZydra/DHBW.Photo-Server/workflows/Go/badge.svg)
# Photo-Server

- [Organisatorisches](#organisatorisches)
- [Aufgabenstellung](#aufgabenstellung)
   - [Anforderungen](#anforderungen)
   - [Optionale Anforderungen](#optionale-anforderungen)
- [WEB-Server in Go](#web-server-in-go)
- [Abgabeumfang](#abgabeumfang)
- [Abgabe](#abgabe)
- [Bewertungskriterien](#bewertungskriterien)

## Organisatorisches
- Die Abgabe des Programmentwurfs erfolgt spätestens am Sonntag, den **10.01.2021** als gepackte Quellen im ZIP-Format per pushci.
- Die Bearbeitung der Aufgabe erfolgt in Gruppen von maximal drei Studierenden.
- Es muss eine Aufstellung der Beiträge der einzelnen Gruppenmitglieder abgegeben
werden.
- Jede Gruppe fertigt eigenständig eine individuelle Lösung der Aufgabenstellung an und reicht diese wie oben angegeben ein.
- Die geforderte Funktionalität muss von Ihnen selbst implementiert werden.
- Sie können sich Anregungen aus anderen Projekten holen, allerdings muss in diesem Fall die Herkunft der Ideen bzw. von Teilen des Quellcodes in den betroffenen Quelldateien vermerkt sein.
- Werden Code-Bestandteile einer anderen Gruppe dieser Veranstaltung übernommen, wird die Abgabe beider Gruppen mit 0% bewertet.
- Leider ist es schon mehrfach vorgekommen, dass Gruppen den Code einer anderen Gruppe übernommen haben, welcher in einem öffentlichen git-Repository gehostet war. Daher muss ich leider dringend empfehlen, für die Zusammenarbeit innerhalb einer Gruppe ein privates Repository zu verwenden.

## Aufgabenstellung
Es soll ein System zur Verwaltung von Photos implementiert werden.

### Eigene Anmerkungen
- Konsole - Für was Datum aus EXIF?

3 Komponenten:  
- Konsolen-Tool
- Webseite
   - Caching
   - BasicAuth
- Backend (REST-API)
   - Jede Anfrage muss im Header Benutzername und Passwort enthalten

**Projektstruktur**  
```
/docs       PDF, UML
/cmd
  /backend
  /website
  /terminal
    main.go
/internal
  # EXIF-Parser
  # Helper
  # ImageManager
  # ...
/web
/test
```

**Speichern der Bilder**  
`Content.csv`enthält pro Bild:  
Zeitstempel | Name | Hashwert (md5)

Für jedes Bild:  
- Bild1.jpeg (Original)
- Previews/Bild1.jpeg (Thumbnail)
- Bild1.json (Hash, Exif, ggf. Kommentare)  
   Bei Kommentaren muss nur der Text ggf. mit Zeitstempel gespeichert werden.

**Benutzerdaten**  
`users.csv`  
Benutzername und Passwort (["salting"](https://play.golang.org/p/tAZtO7L6pm))  
Benutzername darf keine Sonderzeichen enthalten
Benutzername | Passwort

Registieren von Benutzern ergänzen

**Bestellen**  
Die Bilder im Warenkorb werden mit den Metainfos (Auszugformate, Menge) als Zip heruntergeladen.  
Danach kann der Warenkorb gelöscht werden.  
JSON-Datei: Bild mit zugehörigen (Format, Anzal)-Einträgen.

**REST-API**  
- `POST` - Upload eines Fotos:  
  1. Ein Bild pro Aufruf (als Base64)
  2. Bild als Base64 -> Mehrere Bilder, da nur von A-Z, "/" und "+"
  
**Aufgabenverteilung**  
- Benutzerauthentifizierung (Jones)
- Bilder (David)
- Konsole (Patricia)

### Anforderungen
1. Nicht funktional
   1. Die Implementierung hat in der Programmiersprache Go zu erfolgen.
   2. Es dürfen keine Pakete Dritter verwendet werden! Einzige Ausnahme ist ein Paket zur Vereinfachung der Tests.  
      Empfohlen sei hier [www.github.com/stretchr/testify/assert](www.github.com/stretchr/testify/assert)
   3. Alle Source-Dateien und die PDF-Datei müssen die Matrikelnummern aller Gruppenmitglieder enthalten.

2. Allgemein
   1. Die Anwendung soll mehrere Nutzer unterstützen. Die Nutzer sollen komplett voneinander getrennt sein. Es soll keinem Nutzer möglich sein, die Photos anderer Nutzer anzusehen oder abzurufen.
   2. Die Anwendung soll unter Windows und Linux lauffähig sein.
   3. Es soll sowohl Firefox, Chrome und Edge in der jeweils aktuellen Version unterstützt werden. Diese Anforderung ist am einfachsten zu erfüllen, indem Sie auf komplexe JavaScript/CSS ”spielereien“ verzichten. :-)

3. Sicherheit
   1. Die Web-Seite soll nur per HTTPS erreichbar sein.
   2. Der Zugang für die Nutzer soll durch Benutzernamen und Passwort geschützt werden.
   3. Die Passwörter dürfen nicht im Klartext gespeichert werden.
   4. Es soll ”salting“ eingesetzt werden.
   5. Alle Zugangsdaten sind in einer gemeinsamen Datei zu speichern.

4. Photos
   1. Es müssen nur JPEG Dateien unterstützt werden.
   2. Dem Anwender soll eine chronologisch sortierte Liste der Photos angezeigt werden.
   3. Es soll eine Übersicht bestehend aus verkleinerten Vorschaubildern angezeigt
werden. (_Bei Upload verkleinerte Version speichern_)
   4. Es soll möglich sein, ein Photo auch in Originalauflösung anzusehen. (_Click Event -> Popup?_)
   5. Das System soll in der Lage sein, mit tausenden von Photos umzugehen, das Anzeigen einer lange Liste aller Photos ist also nicht geeignet. (_Paging, Dropdown mit verschiedenen Zahlen_)

5. Photo-Upload
   1. Es müssen nur JPEG Dateien unterstützt werden.
   2. Photos sollen über eine WEB-Seite hochgeladen werden können.
   3. Das Datum des Photos soll aus dem Exif-Header der JPEG-Datei entnommen werden. (_Exif-Parser_)
   4. Der mehrfache Upload eines Fotos soll verhindert werden. (_Bei Upload Hash über Pixel erstellen?_)

6. Batch-Upload
   1. Es soll ein Kommandozeilen-Tool zum Upload aller Photos eines Ordners geben.
   2. Als Parameter sollen angegeben werden:
      - Host des Servers
      - Pfad des Ordners mit den Photos
      - Username
      - Passwort
   3. Der Server soll eine geeignete REST-API implementieren, welche vom Kommandozeilen-Tool genutzt werden kann.
   4. Der mehrfache Upload eines Fotos soll verhindert werden.
   5. Auch hier soll das Datum des Photos aus dem Exif-Header der JPEG-Datei entnommen werden.

7. Bestellung von Photos
   1. Aus der Menge aller Photos soll eine Liste mit einer Auswahl an Photos erstellt werden können. (_Auswahlmöglichkeit_)
   2. Jedes Photo dieser Liste soll mit Meta-Informationen versehen werden können, die eine Bestellung von Abzügen ermöglichen.
   3. Dazu gehören mindestens die Anzahl der Abzüge und das Format der Abzüge. (_Standardformate wie 9er ..._)
   4. Es soll möglich sein, einzelne Photos aus dieser Liste zu löschen.
   5. Es soll möglich sein, diese Liste komplett zu löschen.
   6. Es soll möglich sein, die Photos dieser Liste in Originalauflösung zusammen mit den Meta-Informationen als eine ZIP-Datei herunterzuladen.

8. Storage
   1. Alle Informationen sollen zusammen mit den Photos im Dateisystem gespeichert werden. (_Keine DB_)
   2. Es sollen nicht alle Informationen in einer gemeinsamen Datei gespeichert werden.
   3. Es soll ein geeignetes Caching implementiert werden, d.h. es sollen **nicht** bei jedem Request alle Dateien aus dem Dateisystem neu eingelesen werden.
   4. Es soll verhindert werden, dass der Cache über alle Grenzen wächst, wenn sehr viele Fotos verwaltet werden. Gehen Sie davon aus, dass nicht alle Photos in Originalauflösung im Speicher gehalten werden können. (_Server- oder clientseitig?_)

9. Konfiguration
   1. Die Konfiguration soll komplett über Startparameter erfolgen. (Siehe Package `flag`)
   2. Der Port muss sich über ein Flag festlegen lassen.
   3. ”Hart kodierte“ absolute Pfade sind nicht erlaubt.

10. Betrieb
    1. Wird die Anwendung ohne Argumente gestartet, soll ein sinnvoller ”default“ gewählt werden.
    2. Nicht vorhandene aber benötigte Order sollen ggfls. angelegt werden. (_Prüfung und Ordner erstellen beim Start_)
    3. Die Anwendung soll zwar HTTPS und die entsprechenden erforderlichen Zertifikate unterstützen, es kann jedoch davon ausgegangen werden, dass geeignete Zertifikate gestellt werden. Für Ihre Tests können Sie ein ”self signed“ Zertifikat verwenden. Es ist nicht erforderlich zur Laufzeit Zertifikate zu erstellen o.ä.. Ebenso ist keine Let’s Encrypt Anbindung erforderlich.

### Optionale Anforderungen
11. Kommentierung eines Photos
    1. Es soll möglich sein, die Photos zu kommentieren.
    2. Die Kommentare sollen durchsucht werden können.

12. Photogruppen
    1. Photos sollen zu Gruppen zusammengefasst werden können.
    2. Es soll möglich sein, alle Photos dieser Gruppe anzuzeigen. (_Gruppe muss gespeichert werden_)
    3. Eine solche Gruppe soll als Bestellliste übernommen werden können.
    4. Es soll eine ”Dia-Show“ einer Gruppe möglich sein.

## WEB-Server in Go
Es gibt einige interessante Techniken, welche die Implementierung von WEB-Servern in Go unterstützen.  
Der folgende sehr gelungene Vortrag sei empfohlen:  
Mat Ryer: ”Building APIs“  
Golang UK Conference 2015  
https://youtu.be/tIm8UkSf6RA

Eine weitere Quelle für Best-Practices findet sich in diesem Vortrag:  
Peter Bourgon: ”Best Practices in Production Environments“  
QCon London 2016  
https://peter.bourgon.org/go-best-practices-2016/

## Abgabeumfang
- Der komplette Source-Code incl. der Testfälle.
- Dokumentation des Source-Codes  
Die Dokumentation setzt sich aus zwei Bestandteilen zusammen: Zum einen aus der Dokumentation des Sourcecodes. Diese erfolgt am geeignetsten im Sourcecode selbst. Zum anderen aus der Dokumentation der Architektur. Diese soll geeignet sein, einem Außenstehenden einen Überblick über das Gesamtprojekt zu verschaffen, indem Fragen beantwortet werden wie: ”Aus welchen Komponenten besteht das System?“ oder ”Wie arbeiten die Komponenten zusammen?“. Bei Objektorientierten Projekten bietet sich z.B. ein Klassendiagramm an, um die Beziehungen zwischen den Klassen darzustellen, ist allein aber nicht ausreichend.
- Die Abgabe soll eine PDF-Datei beinhalten, welche mindestens die folgenden Kapitel hat:
   - Architekturdokumentation
   - Anwenderdokumentation
   - Dokumentation des Betriebs.
   - Jedes Gruppenmitglied ergänzt eine kurze Beschreibung des eigenen Beitrags zur Projektumsetzung ab (eine Seite reicht).
- Für die Bewertung werden nur die Sourcen und die eine PDF-Datei herangezogen.
- Alle Source-Dateien und die PDF-Datei müssen die Matrikelnummern aller Gruppenmitglieder enthalten.

## Abgabe
Die Abgabe soll per pushci (**siehe PDF**) erfolgen. Hierbei handelt es sich um einen Dienst, welcher aus dem Intranet und dem Internet erreichbar ist.

Dieser Dienst übernimmt Ihre Datei, packt diese aus, überprüft einige formale Kriterien, compiliert die Sourcen und führt alle Tests aus. Eine Datei kann nur abgegeben werden, wenn alle Testfälle ”grün“ sind.

Wer bereits mit einem Continuous-Integration-System gearbeitet hat, ist mit diesem Vorgehen sicher vertraut. Für Studierende, die noch nicht mit einem CI-System gearbei- tet haben, ist das möglicherweise ungewohnt. Es ist daher dringend angeraten, sich mit diesem Dienst frühzeitig vertraut zu machen. Es wäre sehr ungeschickt, erst am Tag der Abgabe um 23:55 das erste mal zu versuchen, eine Datei ”hochzuladen“.

Sie können beliebig oft eine Datei ”hochladen“ und überprüfen lassen. Der Dienst ist während des ganzen Semesters verfügbar. Wenn diese Überprüfung erfolgreich war, kann die Datei abgegeben werden, Sie können die Abgabe jedoch auch bei erfolgreicher Prüfung einfach abbrechen. Bitte geben Sie nur einmal eine finale Lösung ab, indem Sie nach der Überprüfung die Matrikelnummern der Gruppenmitglieder angeben. Nur im Ausnahmefall sollte eine Abgabe mehrfach erfolgen, wobei nur die zuletzt abgegebene Variante bewertet wird.

## Bewertungskriterien
Ihr Programmentwurf wird nach folgenden Kriterien bewertet:
1. Abbildung der Anforderungen (s.o.)
2. Strukturierung und Konsistenz des Quellcodes
3. Ausreichende Testabdeckung des implementierten Codes. (gemessen mit `go test -cover`)
4. Sinnhaftigkeit der Tests (Sonderfälle, Grenzfälle, Orthogonalität usw.)
5. Qualität von Kommentaren und Dokumentation
6. Benutzerfreundlichkeit der Schnittstellen (APIs)
7. Optionale Features
8. Das Design der Web-Seite spielt keine Rolle solange sie benutzbar ist.
