# Wiki Knowledge

This is a wiki article look-up API. It was originally made for a chatbot product that was a real-world business project.

## 1. Format

Single-round task-based chatbot

## 2. Data Source

Baidu Baike:

https://baike.baidu.com/

## 3. Triggers

- **Keyword** wiki
- what is **Keyword**

> Examples:
>
> Beatles wiki
>
> what is Shanghai
>
> ...

## 4. Response

### 4.1 If look-up was successful:

If the look-up was successful, first return the synopsis, then return a message card containing the link to the relevant article.

#### 4.1.1 Message Card

Message card:

- Title: Beatles wiki

- Subtitle: Encyclopedia knowledge helps you know the world better.

- Thumbnail: refer to blueprint

- Link: the link to article Beatles in Baidu Baike

### 4.2 If an error occurred:

If an error occurred during the look-up, give the following reply:

Sorry, there was an error in looking up **Keyword**.

## Details

- When the server is up, there will be two APIs:
  - an API for returning text messages (synopsis)
  - an API for returning message cards
- The chatbot should send a request to the first API and get the text response and card id. From that, it can use the card id to request a message card from the second API.
- Frequency settings:
  - results from crawlers: kept in local files for 1 hour
  - message cards: kept in local files for 1 minute

## Local Test

```bash
./run.sh
```
