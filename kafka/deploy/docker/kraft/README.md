# Kafka KRaft Mode
## kafka
Kafka 설정은 [링크](https://docs.confluent.io/platform/current/installation/docker/config-reference.html)에 들어가면 더 자세히 확인하실 수 있습니다.

### configure
- **KAFKA_NODE_ID : Kafka 노드 ID를 설정합니다. 노드 ID는 각 노드의 식별자 입니다. (Required for KRaft mode)**
- **KAFKA_PROCESS_ROLES : 서버의 역할을 지정합니다. 서버는 controller, broker 혹은 둘의 역할을 다 할 수 있습니다. (Required for KRaft mode)**
- **KAFKA_CONTROLLER_QUORUM_VOTERS : zookeeper.connect 와 대응하는 역할을 하고 있으며, Controller Quorum에 연결할 수 있는 노드를 결정합니다. 식별자는 {id}@{host}:{port} 로 구분을 합니다. (Required for KRaft mode)**
- **KAFKA_CONTROLLER_LISTENER_NAMES : Zookeeper 모드 브로커는 이 값을 설정하면 안되며, Controller에서 사용하는 리스터 이름입니다. (Required for KRaft mode)**
- KAFKA_LISTENER_SECURITY_PROTOCOL_MAP : 리스너 별로 사용할 보안 프로토콜의 키/값 쌍입니다. 해당 코드에서는 모든 프로토콜이 PLAINTEXT와 쌍을 이루도록 설정을 하였습니다.
- KAFKA_ADVERTISED_LISTENERS : Docker 내부에서 사용할 포트와 Docker 네트워크안에서 각 컨테이너에서 연결할 포트를 설정합니다.
- KAFKA_LISTENERS : Kafka가 바인딩하는 {host}:{port} 로 구분을 한 목록입니다.
- KAFKA_INTER_BROKER_LISTENER_NAME : 브로커 간 통신에 사용할 수신 이름을 정의합니다. 서로 통신을 할 때 필요한 매개변수 입니다.
- CLUSTER_ID : 클러스터의 고류 식별자를 지정하는 문자열입니다. 동일한 클러스터에서 동작을 하기 때문에 같은 CLUSTER_ID 로 설정합니다.

## kafka-ui
### configure
- KAFKA_CLUSTERS_0_NAME : Kafka 클러스터 이름을 지정합니다.
- KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS : Kafka 클러스터의 bootstrap 서버를 설정합니다.

## Referenced
[Kafka KRaft 모드](https://medium.com/mo-zza/kafka-kraft-%EB%AA%A8%EB%93%9C-with-docker-%EB%8F%99%EB%AC%BC%EC%9B%90%EC%9D%84-%ED%83%88%EC%B6%9C%ED%95%9C-kafka-8b5e7c7632fa)