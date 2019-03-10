const fs = require('fs');
const path = require('path');
const yaml = require('js-yaml');

const config = yaml.safeLoad(fs.readFileSync(path.resolve(__dirname, './config.yaml'), 'utf8'));
// console.log('aws:SourceIp', config.Authorization['whitelist_ip']);
function generatePolicy(principal_id, effect, resource, user=null) {
  return {
    principalId: principal_id,
    policyDocument: {
      Version: '2012-10-17',
      Statement:
        [
          {
            Action: 'execute-api:Invoke',
            Effect: effect,
            Condition: {
              IpAddress: {
                'aws:SourceIp': config.Authorization['whitelist_ip'],
              }
            },
            Resource: resource
          }
        ]
    },
    context: user
  };
}


exports.handler = async (event, context) => {
  try {
    return generatePolicy('id', 'Allow', event.methodArn, {user: 'test'});
  }catch (error) {
    console.error(error);
    throw error;
  }

}