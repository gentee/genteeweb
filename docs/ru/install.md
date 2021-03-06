# Установка веб-сервера

## Шаг 1. Скачайте программу

Скачайте версию веб-сервера для вашей операционной системы. Если у вас установлен язык программирования [Go](https://golang.org/), то вы можете сами скомпилировать исполняемый файл.


* [Linux 64-bit](/)
* [Windows 64-bit](/)
* [macOS 64-bit](/)

## Шаг 2. Запуск локального сервера

Мы рекомендуем в начале запустить _GenteeWeb_ в качестве веб-сервера на локальном компьютере. Если вы планируйте использовать _GenteeWeb_ только как генератор статических страниц, то вам будет достаточно этого шага. Перед запуском нужно создать [файл конфигурации](config.html) _genteeweb.yaml_ и папку со [страницами и шаблонами](content.html).  
Если вы хотите быстро посмотреть на работу программы, то скачайте и распакуйте демонстрационный .zip архив и при запуске _genteeweb_ укажите путь к файлу _genteeweb.yaml_. Например, вы распаковали demo.zip в поддиректорию _c:\temp\demo_. В этом случае, запустите веб-сервер командой

```bash
genteeweb.exe "c:\temp\demo\examples\genteeweb.yaml"
```

Для Linux это может выглядеть так
```bash
./genteeweb "/home/ak/temp/demo/examples/genteeweb.yaml"
```

Если вы скопируйте запускаемый файл _genteeweb_ в директорию, где находится файл конфигурации, то в этом случае веб-сервер можно запускать без параметров.

## Шаг 3. Запуск сервера на хостинге.

