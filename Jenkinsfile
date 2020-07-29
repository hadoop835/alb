library "alauda-cicd"
def language = "golang"
AlaudaPipeline {
    config = [
        agent: 'golang-and-devops',
        folder: 'src/alb2',
        chart: [
            [
                enabled: false
            ]
        ],
        scm: [
            credentials: 'acp-acp-bitbucket-new'
        ],
        docker: [
            credentials: "alaudak8s",
            enabled: false,
        ],
        sonar: [
            binding: "sonarqube",
            enabled: false,
        ],
    ]
    env = [
        GO111MODULE: "on",
        GOPROXY: "https://goproxy.cn,direct",
    ]
    steps = [
    ]
    yaml = "alauda.yaml"
}
