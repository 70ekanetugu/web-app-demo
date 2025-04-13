# About
簡易的なTODOアプリケーション用WebAPI。AWSなどで検証する際に使用する。  

# 事前準備
GitHub ActionsでECRに登録できるようになっていますが、その前に以下の操作が必要。 

## AWS IAM IDプロバイダの登録
GitHub ActionsからAWSリソースにアクセスするための一時的な認証情報を得ることができるようにGitHubをOIDCプロバイダーとして登録する。  
1. AWSコンソールにログインし、IAMに遷移する
2. IAMのIDプロバイダを選択し、プロバイダ追加を押す
3. 以下の通り選択・入力し、プロバイダを追加する
   - プロバイダのタイプ： OpenID Connect
   - プロバイダのURL: https://token.actions.githubusercontent.com
   - 対象者： sts.amazonaws.com
   - タグ： 任意

## IAM Roleの作成
GitHub Actionsで使用するためのIAM Roleを作成する。  
1. AWSコンソールにログインし、IAMに遷移する
2. IAMのロールを選択し、以下のロールを作成する。
   - 信頼ポリシー
     ```json
     {
      "Version": "2012-10-17",
      "Statement": [
         {
            "Effect": "Allow",
            "Principal": {
               "Federated": "arn:aws:iam::<アカウントID>:oidc-provider/token.actions.githubusercontent.com"
            },
            "Action": "sts:AssumeRoleWithWebIdentity",
            "Condition": {
               "StringEquals": {
                  "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
               },
               "StringLike": {
                  "token.actions.githubusercontent.com:sub": "repo:<リポジトリオーナー名>/<リポジトリ名>:*"
               }
            }
        }]
     }
     ```
   - 許可ポリシー
     ```json
     {
         "Version": "2012-10-17",
         "Statement": [
         {
               "Sid": "GetAuthorizationToken",
               "Effect": "Allow",
               "Action": [
                  "ecr:GetAuthorizationToken"
               ],
               "Resource": "*" // GetAuthorizationTokenは制限ができないため「*」指定。
         },
         {
               "Sid": "PushImageOnly",
               "Effect": "Allow",
               "Action": [
                  "ecr:BatchCheckLayerAvailability",
                  "ecr:InitiateLayerUpload",
                  "ecr:UploadLayerPart",
                  "ecr:CompleteLayerUpload",
                  "ecr:PutImage",
                  "ecr:BatchGetImage" // Manifestの取得に必要
               ],
               "Resource": "arn:aws:ecr:ap-northeast-1:<アカウントID>:repository/*"
         }]
     }
     ```

## ECRでリポジトリを作成
イメージを登録する用のリポジトリを作成しておく。  
必要なリポジトリは以下２つ。  
- todo-go-app
- fluentbit

## GitHub Actions用のシークレット登録
以下のシークレットを登録しておく。  
- AWS_ACCOUNT_ID： AWSのアカウントID
- AWS_ROLE_FOR_GITHUB_ACTIONS： 前述で作成したIAMロールのarn
- BUCKET_NAME： goアプリのバイナリを格納するバケットの名前

# Usage
## 実行
ローカル上で実行する場合は以下の通り。  
```shell
docker compose up -d
```

- [ ] http://localhost/hello
- [ ] http://localhost/todos

## バイナリビルド手順
コンテナに潜ってビルドしたい時は以下の通り。

1. ビルド
   ```shell
   go build -o ./server .
   ```

## ECRへの登録手順
### GitHub Actionsを使う場合
1. GitHub Actionsに遷移し、 `Push Docker Image to ECR` で `Run workflow` を実行する
2. 以下に登録できていることを確認する
   - [ ] ECRの todo-go-app
   - [ ] ECRの fluentbit
   - [ ] S3の指定バケット

### 手動でAWS cliを使う場合
手動でaws cliを使ってECRへイメージを登録する手順は以下の通り。  

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

# その他
## Fluent Bit
### 基礎
docker-composeのログドライバオプションの `fluentd-address` は、ホストのアドレス:ポートを見に行くため、127.0.0.1を指定する必要がある。  
(サービス名fluent-bitは不可)  
また、ログドライバを変更する場合、先にFluent-Bitが起動している必要があるため、web, go-appのdepend_onにfluent-bitを指定している。  
このため、fluent-bit側でヘルチェックが必要なため、docker/fluentbit/fluent-bit.conf ではヘルスチェック用にHTTP_ServerをOnにしている。

### 設定
docker/fluentbit/fluent-bit.conf で設定を修正可能。  
詳細は公式を参照。  
https://docs.fluentbit.io/manual/pipeline/pipeline-monitoring


