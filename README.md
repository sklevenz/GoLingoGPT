# GoLingoGPT

A server that corrects grammar with help of ChatGPT. 

## try it out

```
export OPENAI_KEY="your openai api key"
#export OPENAI_MOCK="true"  // use a mock instead of calling openai

go run golingogpt_server.go

./bin/testCorrectGrammerGet.sh
./bin/testCorrectGrammerPost.sh
```
