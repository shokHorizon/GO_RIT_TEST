# Гайд по использованию исполнителя блок-схем
## Подготовка рабочей среды
1) Установите Golang 1.20
2) Скачайте необходимые библиотеки
## Описание скрипта
Скрипт описывается в json формате и состоит из двух блоков:

1) Actions - Обычные действия, по окончании которых запускаются следующие
2) Conditions - Условия, которые сравнивают значения и исполняют следующие блоки, исходя из результата сравнения

### Actions
    ```json
    {
        "name": "appendTest1",    # Имя блока (должно быть уникальным)
        "action": "appendString", # Тип действия
        "params": {               # Параметры действия (если не указаны
                                    необходимые - возьмет из предыдущих
                                    блоков, иначе - ошибка)
            "file": "sample.txt",
            "text":"!!test1\n"
        },
        "next": ["renameFile"]    # Действия или условия, следующие за текущим
    },
    ```

#### Типы действий:
    
1) timeFromString - Преобразует время из строки, чтобы передать его блоку сравнения. 
    Обязательный параметр - 'time', где указывается время в формате 
    ЧЧ:ММ ДД.ММ.ГГГГ (Пример: 12:30 01.01.2022)

2) сreateFile - Создает файл с указанным названием
    Обязательный параметр - 'file', содержащий название файла

3) кenameFile - Переименовывает файл
    Обязательный параметр - 'file', содержащий текущее название файла
                          - 'rename', содержащий новое название

4) appendString - Добавляет текст в файл
    Обязательный параметр - 'file', содержащий название файла
                          - 'text', содержащий текст к добавлению

5) getCreationTime - Получает время создания файла и передает в блок сравнения
    Обязательный параметр - 'file', содержащий название файла

### Conditions
    ```json
    {
        "name": "creationTimeCondition",        # Имя блока (уникальное)
        "action": "ifTime",                     # Тип сравнения
        "params": {                             # Параметры действия (если не
                                                  указаны необходимые -
                                                  возьмет из предыдущих
                                                  блоков, иначе - ошибка)
            "first_arg": "fileCreationTime",    # Названия первого блока для
                                                  сравнения
            "second_arg": "fixedCreationTime",  # Названия Второго блока для
                                                  сравнения
            "operator": "<"                     # Оператор сравнения
        },
        "next": ["appendTest1", "appendTest2"]  # Первое действие исполнится
                                                  при положительном результате,
                                                  второе - при отрицательном
    }
    ```

### Пример скрипта назван "example.json"

## Использование
> go run main.go "Название файла со скриптом"

Пример: 
> go run main.go example.json

Результат исполнения можно увидеть в файле "Название вашего файла".log