# About
簡易的なTODOアプリケーション用WebAPI。AWSなどで検証する際に使用する。  

# Usage
## バイナリ生成
1. 以下コマンドを実行するだけ。
   ```shell
   go build -o ./server .
   ```

## イメージ作成・登録
ECSで使用する場合は以下手順でECRへイメージを登録する。
1. アカウントIDの取得
   ```shell
   AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text [--profile プロファイル名])
   ```
2. ECRに登録するイメージ(タグ)の作成
   ```shell
   docker build -t todo-go-app -f ./docker/golang/Dockerfile .
   docker image tag todo-go-app:latest ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/todo-go-app:v1
   ```
3. ECRに `todo-go-app` という名前のリポジトリを作成しておく 
4. 認証
   ```shell
   aws ecr --region ap-northeast-1 get-login-password [--profile プロファイル名] | docker login --username AWS --password-stdin https://${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/todo-go-app:v1
   ```
4. ECRにpushする
   ```shell
   docker image push ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/todo-go-app:v1
   ```
