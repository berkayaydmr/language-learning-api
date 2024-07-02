# Language Learning API

## Project Layout

```
- cmd
    - service
        - main.go
- pkg
    - transport
        - transport.go
        - transport_test.go
    - storage
        - storage.go
        - storage_test.go
- README.md
```

## Endpoints

### List

- **URL:** `/words`
- **Method:** `POST`
- **Headers:**
  - **Content-Type:** `application/json`
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Content:**
    ```json
    [
        {
            "id": 1,
            "word": "apple",
            "translation": "elma",
            "language": "English",
            "exampleSentence": "I eat an apple every day."
        },
        {
            "id": 2,
            "word": "banana",
            "translation": "muz",
            "language": "English",
            "exampleSentence": "Bananas are rich in potassium."
        }
    ]
    ```

- **Error Responses:**
  - **Not Found Error:**
    - **Status Code:** `404 Not Found`
    - **Content:**
      ```json
      {
          "error": "No word found."
      }
      ```
  - **Internal Server Error:**
    - **Status Code:** `500 Internal Server Error`
    - **Content:**
      ```json
      {
          "error": "Failed to list words."
      }
      ```
