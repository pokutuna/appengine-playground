error-reporting-go
===

## GAE Standard Environment

- panic でアプリケーションコンテナ自体が異常終了したら拾ってくれる
- それ以外で Stderr に書いても Reporting は拾ってくれない(Logging は拾ってくれる)
- Stackdriver Error Reporting クライアントを利用すると拾ってくれる
  - API 有効にしていない場合は以下のエラー
    ```
    code = PermissionDenied desc = Stackdriver Error Reporting API has not been used in project <project_id> before or it is disabled. Enable it by visiting https://console.developers.google.com/apis/api/clouderrorreporting.googleapis.com/overview?project=<project_id> then retry. If you enabled this API recently, wait a few minutes for the action to propagate to our systems and retry.
    ```
  - Stackdriver Reporting API を有効にすると GAE インスタンスが再起動される?
    - `Quitting on terminated signal`
    - めちゃトラフィック来てたら怖い気がするが...
  - client 使うと普通に ErrorReporting に記録される
  - `errorreporting.Config` に `ServiceVersion` 渡さないとバージョン情報が入らない
    ```go
    errorreporting.Config{
      ServiceName: "error-reporting-go",
    }
    ```
  - ソースコードへの参照は機能するが、直前の Stackdriver Debugging で見てるサービスが一致していないとソースツリーが意図しないものになりそう

## GAE Flex Environment
