properties([
  parameters([
    [
      $class: 'CascadeChoiceParameter',
      choiceType: 'PT_SINGLE_SELECT',
      name: 'Account',
      description: 'Select the account',
      script: [
        $class: 'GroovyScript',
        fallbackScript: [
          classpath: [],
          sandbox: true,
          script: 'return ["ERROR: Failed to fetch accounts"]'
        ],
        script: [
          classpath: [],
          sandbox: true,
          script: '''
            import groovy.json.JsonSlurper
            def url = new URL('https://662cab7c5e116819738b01fe:supertoaster@toaster.altuhov.su/api/dimension/demo-org/account?workspace=master&fallbacktomaster=true&needdimdata=false')
            def connection = url.openConnection()
            connection.setRequestMethod('GET')
            
            try {
                def response = new JsonSlurper().parse(connection.inputStream)
                if (response.Dimensions) {
                    return response.Dimensions.collect { it.DimValue }.unique()
                } else {
                    return ["No accounts found"]
                }
            } catch (Exception e) {
                return ["Error: ${e.message}"]
            }
          '''
        ]
      ]
    ]
  ])
])

pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        echo "Selected account: ${params.Account}"
      }
    }
  }
}
