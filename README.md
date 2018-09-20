# eosair

## About eosair

eosair는 EOS 계정을 수집하고 수집된 계정에 memo 광고(0.0001 EOS를 보내면서 memo 에 광고 메시지를 담는 것)를 보내거나 EOS 토큰을 에어드랍하는 기능을 갖추고 있습니다.

## 사용 전 필수 설치 프로그램

eosair를 사용하려면 docker와 eosio 컨테이너가 실행된 상태여야 합니다.

## docker 설치

https://docs.docker.com/

## docker 에서 eosio 컨테이너 실행
docker 컨테이너 실행 후 wallet plugin 이 있어야 하고 wallet 이 존재해야합니다. memo 광고를 하려면 default wallet 으로 사용할 private key를 import 해주세요. EOS 토큰을 에어드랍하려면 EOS 토큰이 default wallet의 계정에 충분히 들어있어야 합니다. 실행이 되지 않는 경우 계정의 Ram, CPU, NET 자원 정보를 확인하세요.
```{r, engine='bash'}
$ docker pull eosio/eos
$ sudo docker run -it --name eosio -p 8888:8888 -p 9876:9876 eosio/eos /bin/bash
$ docker start eosio
```

## EOS 계정 수집 방법
EOS 계정 수집은 ./eosair 기본 명령어로 실행할 수 있습니다. 예를 들어, 980만 번째 블록부터 최신블록까지 계정을 찾는 경우 -n 9800000 을 함께 입력하면 됩니다.  전체 계정을 수집하는 경우 수십일이 걸릴 수도 있습니다. 수집된 계정은 홈 디렉토리에 accounts.csv 파일로 저장됩니다. -n 옵션을 입력하지 않으면 accounts.csv 파일의 가장 최신 블록 숫자를 기준으로 계정 수집을 재개합니다.
```{r, engine='bash'}
$ ./eosair -n 9800000
```
## EOS memo 광고 실행
예시) accounts.csv 에 계정 정보가 수집된 상태이고 3번째 행의 계정부터 1만명의 계정한테 memo 광고를 전달하려고 한다.
```{r, engine='bash'}
$ ./eosair -memo "메모 내용" -sender [EOS 계정명] -n 3 -p [EOS default wallet unlock 패스워드]
```
## EOS 에어드랍 실행
위의 memo 광고처럼 accounts.csv 에 계정 정보가 미리 저장되어 있어야 한다. 예를 들어, 3번째 행부터 실행하려면 -n 3 옵션을 입력한다. 에어드랍은 유저가 EOS를 보유한 비율대로 지급할 수도 있고 동일한 수량을 똑같이 지급할 수도 있다.
```{r, engine='bash'}
$ ./eosair -airdrop [토큰 심볼 이름] -contract [컨트랙트 계정이름] -sender [EOS 계정이름] -amount "1.0000/EOS" -memo "메모 내용" -n 3 -p [default wallet unlock 패스워드]
```
## Donation

BTC : 1JXFbzhtr1rFVrrAmyRgvEZLdZapAgczHh<br>
EOS : givemeeosplz

## Contact Information

* [Twitter](https://twitter.com/booyoun)

## LICENSE

MIT License Copyright © 2018 Sushimania
