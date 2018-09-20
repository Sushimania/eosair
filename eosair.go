package main

import (
	"fmt"
	"os"
	"os/exec"
	"encoding/json"
	"strconv"

	"eosair/model"
	"encoding/csv"
	"bufio"
	"flag"
	"sync"
	"time"
	"os/user"
	"eosair/util"
	"strings"
)

var (
	accounts []string
	wr *csv.Writer
	// 고루틴 싱크 라이브러리 (고루틴 제어)
	wait sync.WaitGroup
	programState string

	saveFileReadComplete bool
	savedEndBlockNum int
	dir string
	savedBlock []string
	savedAccount []string

	// 커맨드라인 매개변수
	startBlockNum int
	memo string
	sender string
	unlockPassword string
	tokenSymbol string
	contract string
	amount string
)

func main() {
	// docker 에 eosio 컨테이너가 있는 상태여야함(start 된 상태여야함)

	// 커맨드 매개변수 입력받기
	arg1 := flag.Int("n", 1, "처음 EOS 계정찾기를 시작할 때 몇 번째 블록부터 찾을지 설정할 수 있습니다. 단, ...userHomePath/accounts.csv 가 존재하면 마지막으로 기록된 블록 숫자보다 클 때만 유효합니다.\n")
	arg2 := flag.String("memo", "", "...userHomePath/accounts.csv 파일의 계정에 0.0001 EOS 를 memo 와 함께 보냅니다. memo 에 들어갈 텍스트를 함께 입력하세요. (256 bytes 이하. default wallet 이 존재해야함)\n")
	arg3 := flag.String("sender", "", "EOS 또는 토큰 전송시 보내는 계정 이름을 입력하세요. docker 의 eosio 컨테이너 안에 생성한 default wallet 에 import 된 계정이어야 합니다.\n")
	arg4 := flag.String("p", "", "docker 의 wallet unlock 을 위한 password 를 입력하세요.\n")
	arg5 := flag.String("airdrop", "", "에어드랍하려는 토큰 심볼 이름을 입력하세요.\n")
	arg6 := flag.String("contract", "", "에어드랍하려는 토큰의 스마트컨트랙트 계정을 입력하세요.\n")
	arg7 := flag.String("amount", "1.0000/EOS", "에어드랍하려는 토큰 수량을 입력하세요. -amount \"2.0000/EOS\" 를 입력하면 20.0000 EOS 를 가진 계정에는 40.0000 토큰이 지급됩니다. -amount \"1.0000\"과 -amount \"1.0000/EOS\"가 서로 다른 의미라는걸 기억하세요. -amount \"1.0000\"는 각 계정에 1.0000 토큰을 지급한다는 의미이고 -amount \"1.0000/EOS\"는 각 계정의 EOS 잔액을 확인하고 EOS 수량만큼 1:1 비율로 토큰을 지급한다는 의미입니다.\n")
	flag.String("v", "", "v1.0.10\nCopyright (c) 2018, Booyoun Kim\nEOS 계정 수집, memo 광고, 에어드랍이 가능한 프로그램입니다.(docker 가 설치되어 있어야함)\n\n1. docker 에서 eosio 이미지 가져오는 방법\n$ docker pull eosio/eos\n\n2. 컨테이너 생성\n$ sudo docker run -it --name eosio -p 8888:8888 -p 9876:9876 eosio/eos /bin/bash\n\n3. 컨테이너 시작\n$ docker start eosio\n\n4. EOS 계정 수집(980만 번째 블록부터 최신블록까지 계정을 찾는 경우 - 전체 계정을 수집하는 경우 수십일이 걸릴 수도 있습니다.)\n$ ./eosair -n 9800000\n\n5. memo 대량 전송(accounts.csv 에 계정 정보가 저장되어 있어야한다. 3번째 행부터 실행하려면 -n 3 옵션을 입력한다.)\n$ ./eosair -memo \"메모내용\" -sender [EOS 계정이름] -n 3 -p [default wallet unlock 패스워드]\n\n6. 에어드랍(accounts.csv 에 계정 정보가 저장되어 있어야 한다. 예를 들어, 3번째 행부터 실행하려면 -n 3 옵션을 입력한다.\n$ ./eosair -airdrop [심볼 이름] -contract [컨트랙트 계정이름] -sender [EOS 계정이름] -amount \"1.0000/EOS\" -memo \"메모 내용\" -n 3 -p [default wallet unlock 패스워드]\n")
	flag.Parse()

	startBlockNum = *arg1
	memo = *arg2
	sender = *arg3
	unlockPassword = *arg4
	tokenSymbol = *arg5
	contract = *arg6
	amount = *arg7

	if startBlockNum < 1 {
		startBlockNum = 1
	}

	// ------------------------------------- accounts.csv 파일 읽기 START -------------------------------------
	savedEndBlockNum = -1
	// csv 파일 저장 위치
	u, _ := user.Current()
	dir = u.HomeDir
	saveFileReadComplete = false

	// 배열에 저장 (숫자, 문자로 구성)
	file1, err := os.Open(dir + "/accounts.csv")
	if err != nil {
		// accounts.csv 파일이 존재하지 않음
		fmt.Printf("EOS 계정 수집을 처음부터 시작합니다.\n")
	} else {
		// csv reader 생성
		rdr := csv.NewReader(bufio.NewReader(file1))

		// csv 내용 모두 읽기
		rows, _ := rdr.ReadAll()

		// 행,열 읽기
		for i := range rows {
			// 기존 accounts.csv 파일에 저장된 계정 정보
			accounts = append(accounts, rows[i][0] + "," + rows[i][1])
			// 마지막 블록 숫자 확인용
			savedEndBlockNum, _ = strconv.Atoi(rows[i][0])

			savedBlock = append(savedBlock, rows[i][0])
			savedAccount = append(savedAccount, rows[i][1])
		}

		saveFileReadComplete = true
	}
	// ------------------------------------- accounts.csv 파일 읽기 END -------------------------------------

	if tokenSymbol != "" {
		programState = "airdrop"
	} else {
		if memo != "" {
			programState = "sendMemo"
		} else {
			programState = "collectAccount"
		}
	}

	switch programState {
		case "collectAccount":
			collectAccount()
			break
		case "sendMemo":
			sendMemo()
			break
		case "airdrop":
			airdrop()
			break
	}
}

func airdrop() {
	if tokenSymbol == "" {
		fmt.Printf("-airdrop 옵션으로 토큰 심볼 이름을 함께 입력하세요.\n")
		return
	}

	if contract == "" {
		fmt.Printf("-contract 옵션으로 토큰의 컨트랙트 계정 이름을 함께 입력하세요.\n")
		return
	}

	if sender == "" {
		fmt.Printf("-sender 옵션으로 sender 계정 이름을 함께 입력하세요.\n")
		return
	}

	if unlockPassword == "" {
		fmt.Printf("-p 옵션으로 wallet unlock 패스워드를 함께 입력하세요.\n")
		return
	}

	// 시작 번호를 별도로 지정하는 경우
	startNum := startBlockNum - 1
	lengthAccounts := len(savedAccount)
	if lengthAccounts - startNum <= 0 {
		fmt.Printf("토큰을 보낼 계정이 더 이상 없습니다. (-n 숫자는 accounts.csv 파일의 레코드 개수보다 작아야함)\n")
		return
	}

	fmt.Printf("타겟 계정 수 : %v\n", lengthAccounts - startNum)
	fmt.Printf("%v 토큰 에어드랍을 시작합니다.\n", tokenSymbol)

	// 토큰 분배시 EOS 잔액 비율에 따라 지급하는 것인가?
	autoRate := false

	// amount 에 /EOS 가 포함되어 있나?
	if strings.Contains(amount, "/EOS") {
		autoRate = true
		// /EOS 문자를 제거
		amount = strings.Replace(amount, "/EOS", "", -1)
	}

	for i := startNum; i < lengthAccounts; i++ {
		// 각 계정으로 토큰 전송
		// BP url 블록 순서대로 선택
		bpUrl := util.SelectBPUrl(i)

		if savedAccount[i] == "" || len(savedAccount[i]) < 12 || strings.Contains(savedAccount[i], ".") {
			continue
		}

		// 토큰 받을 계정의 현재 EOS 잔액 구하기
		if autoRate {
			 // amount 비율 계산 시작
			rateAmount, _ := strconv.ParseFloat(amount, 64)

			// cleos --url https://eos.greymass.com:443 get account givemeeosplz -j
			command := `docker exec eosio cleos --url ` + bpUrl + ` get account ` + savedAccount[i] + ` -j`
			cmd2 := exec.Command("sh", "-c", command)
			output, err := cmd2.CombinedOutput()
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
			jsonBody := []byte(output)
			var accountInfo model.AccountInfo
			_ = json.Unmarshal(jsonBody, &accountInfo)

			unstaked := strings.Replace(accountInfo.CoreLiquidBalance, " EOS", "", -1)
			staked := accountInfo.VoterInfo.Staked
			refundCpu := strings.Replace(accountInfo.RefundRequest.CPUAmount, " EOS", "", -1)
			refundNet := strings.Replace(accountInfo.RefundRequest.NetAmount, " EOS", "", -1)

			unstakedFloat, _ := strconv.ParseFloat(unstaked, 64)
			stakedFloat := float64(staked) / 10000

			// Refund 합 계산
			refundCpuFloat, _ := strconv.ParseFloat(refundCpu, 64)
			refundNetFloat, _ := strconv.ParseFloat(refundNet, 64)
			refund := refundCpuFloat + refundNetFloat

			// 총 잔액
			totalBalance := unstakedFloat + stakedFloat + refund

			// 토큰 지급 비율 계산
			rateAmount = rateAmount * totalBalance
			// 최종 지급할 토큰 수량
			amount = strconv.FormatFloat(rateAmount, 'f', 4, 64)
		}

		if amount == "0.0000" || amount == "0" || amount == "" {
			continue
		}

	tokenTransfer:

		// cleos --url https://api.main.alohaeos.com:443 push action ethsidechain transfer '[ "g44tsojxguge", "moneymakings", "50.0000 EETH", "EETH send"]' -p g44tsojxguge
		command := `docker exec eosio cleos --url ` + bpUrl + ` push action ` + contract + ` transfer '[ "` + sender + `", "` + savedAccount[i] + `", "` + amount + ` ` + tokenSymbol + `", "` + memo + `"]' -p ` + sender
		cmd1 := exec.Command("sh", "-c", command)
		_, err := cmd1.CombinedOutput()
		if err != nil {
			fmt.Printf("transfer 명령을 실행할 수 없는 상태입니다. (error: %v)\n", err)

			// 15분마다 wallet unlock 해줘야함
			unlockWallet(bpUrl, unlockPassword)
			goto tokenTransfer

		} else {
			fmt.Printf("%v : 토큰 에어드랍 완료\n", "(" + strconv.Itoa(i + 1) + "/" + strconv.Itoa(lengthAccounts) + ") " + savedBlock[i] + " block, " + savedAccount[i])
		}
	}

	fmt.Println("All done\n")
}

func sendMemo() {
	if sender == "" {
		fmt.Printf("-sender 를 입력하세요.\n")
		return
	}

	if unlockPassword == "" {
		fmt.Printf("-p 옵션으로 wallet unlock 패스워드를 함께 입력하세요.\n")
		return
	}

	// 시작 번호를 별도로 지정하는 경우
	startNum := startBlockNum - 1
	lengthAccounts := len(savedAccount)
	if lengthAccounts - startNum <= 0 {
		fmt.Printf("memo 를 보낼 계정이 더 이상 없습니다. (-n 숫자는 accounts.csv 파일의 레코드 개수보다 작아야함)\n")
		return
	}

	fmt.Printf("타겟 계정 수 : %v\n", lengthAccounts - startNum)
	fmt.Printf("memo 메시지 전송을 시작합니다.\n")

	for i := startNum; i < lengthAccounts; i++ {
		// 각 계정으로 0.0001 EOS 전송 + memo
		// BP url 블록 순서대로 선택
		bpUrl := util.SelectBPUrl(i)

		if savedAccount[i] == "" || len(savedAccount[i]) < 12 || strings.Contains(savedAccount[i], ".") {
			continue
		}

	eosTransfer:

		// cleos --url https://api.main.alohaeos.com:443 transfer g44tsojxguge yellowhammer "0.0001 EOS" "Send memo text https://www.wannabit.io"
		command := `docker exec eosio cleos --url ` + bpUrl + ` transfer ` + sender + ` ` + savedAccount[i] + ` "0.0001 EOS" "` + memo + `"`
		cmd1 := exec.Command("sh", "-c", command)
		_, err := cmd1.CombinedOutput()
		if err != nil {
			fmt.Printf("transfer 명령을 실행할 수 없는 상태입니다. (error: %v)\n", err)

			// 15분마다 wallet unlock 해줘야함
			unlockWallet(bpUrl, unlockPassword)
			goto eosTransfer

		} else {
			fmt.Printf("%v : memo 전송 완료\n", "(" + strconv.Itoa(i + 1) + "/" + strconv.Itoa(lengthAccounts) + ") " + savedBlock[i] + " block, " + savedAccount[i])
		}
	}

	fmt.Println("All done\n")
}

func collectAccount() {
	// 최신 블록 검색
	fmt.Printf("최신 블록 검색중\n")
	command1 := `docker exec eosio cleos --url https://eos.greymass.com:443 get info`
	cmd1 := exec.Command("sh", "-c", command1)
	output, err := cmd1.CombinedOutput()
	if err != nil {
		fmt.Printf("docker 에서 eosio 컨테이너를 찾을 수 없습니다. (error: %v)\n", err)
		return
	}

	jsonBody := []byte(output)
	var getInfo model.GetInfo
	_ = json.Unmarshal(jsonBody, &getInfo)
	endBlock := getInfo.HeadBlockNum

	// TODO : 테스트
	//endBlock = 9692420
	fmt.Printf("head_block_num: %v\n", endBlock)

	var startBlock int

	// 블록 탐색 시작점
	if savedEndBlockNum == -1 {
		// 첫 계정 탐색
		startBlock = 1
	} else {
		// 기존 저장된게 있다면
		startBlock = savedEndBlockNum + 1
	}

	if startBlock < startBlockNum {
		// accounts.csv 에 저장된 마지막 블록 숫자보다 높은 숫자를 옵션으로 입력하면 그거부터 실행
		startBlock = startBlockNum
	}

	// CSV 하나의 레코드마다 저장방식
	file2, err := os.Create(dir + "/accounts.csv")
	if err != nil {
		panic(err)
	}

	// csv writer 생성
	wr = csv.NewWriter(bufio.NewWriter(file2))

	if saveFileReadComplete {
		// 기존 저장된 내용을 다시 write
		for i := 0; i < len(savedBlock); i++ {
			wr.Write([]string{savedBlock[i], savedAccount[i]})
		}
		wr.Flush()
	}

	count := 0
	for {
		time.Sleep(250 * time.Millisecond)

		if startBlock + count < endBlock {
			wait.Add(1)
			go search(startBlock + count)
			count++
		} else {
			fmt.Printf("계정 수집을 종료합니다.")
			break
		}
	}

	wait.Wait()
	fmt.Println("All done\n")
}

func unlockWallet(bpUrl string, unlockPassword string) {
	// cleos wallet unlock
	command := `docker exec eosio cleos wallet unlock --password ` + unlockPassword
	cmd1 := exec.Command("sh", "-c", command)
	_, err := cmd1.CombinedOutput()
	if err != nil {
		fmt.Printf("wallet unlock 실패\n")
	} else {
		fmt.Printf("wallet unlock 성공\n")
	}
}

func search(n int) {
	searchBlock := n
	fmt.Printf("%v 번째 블록에서 계정 찾는중\n", searchBlock)

	// BP url 블록 순서대로 선택
	bpUrl := util.SelectBPUrl(n)

	command := `docker exec eosio cleos --url ` + bpUrl + ` get block ` + strconv.Itoa(searchBlock)
	cmd1 := exec.Command("sh", "-c", command)
	output, err := cmd1.CombinedOutput()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	jsonBody := []byte(output)
	var block model.Block
	_ = json.Unmarshal(jsonBody, &block)

	// 트랜잭션이 존재하는가?
	length := len(block.Transactions)
	if length > 0 {
		// fmt.Printf("length: %v\n", length)

		for j := 0; j < length; j++ {
			// Trx 안이 어떤 경우에는 string, 어떤 경우에는 json 이라서 Trx 가 string 인지 먼저 확인
			// fmt.Printf("block.Transactions[" + strconv.Itoa(j) + "].Trx: %v", block.Transactions[j].Trx)

			if block.Transactions[j].Trx != "" {
				// Trx 가 json 이 아닌 값
				continue
			}

			var blockDetail model.BlockDetail
			_ = json.Unmarshal(jsonBody, &blockDetail)

			// 혹시 BP 관련 트랜잭션인가?
			if blockDetail.Transactions[j].Trx.Transaction.Actions[0].Account != "eosio.token" {
				continue
			}

			tempAccount := blockDetail.Transactions[j].Trx.Transaction.Actions[0].Data.From
			accounts = append(accounts, strconv.Itoa(searchBlock)+","+tempAccount)

			// csv 내용 쓰기
			wr.Write([]string{strconv.Itoa(searchBlock), tempAccount})

			tempAccount = blockDetail.Transactions[j].Trx.Transaction.Actions[0].Data.To
			accounts = append(accounts, strconv.Itoa(searchBlock)+","+tempAccount)

			// csv 내용 쓰기
			wr.Write([]string{strconv.Itoa(searchBlock), tempAccount})
			wr.Flush()
		}
	}
	wait.Done()
	fmt.Printf("%v 번째 블록 탐색 종료\n", n)
}