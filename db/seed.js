// const uuid = require('node-uuid');
const AWS = require('aws-sdk');
const yaml = require('js-yaml');
const fs = require('fs');
const path = require('path');
const getUuid = require('uuid-by-string')
var log4js = require('log4js');
var logger = log4js.getLogger();
logger.level = 'debug';


const configPath = path.resolve(__dirname, '../config.yml');
const config = yaml.safeLoad(fs.readFileSync(configPath, 'utf8'));

AWS.config.update({ region: config.region });
const dynamodb = new AWS.DynamoDB.DocumentClient();

const serviceTableName = `${config.service}-${config.stage}-swagger-dynamo-serviceinfo`;
const versionTableName = `${config.service}-${config.stage}-swagger-dynamo-versioninfo`;


const serviceInfo = [{
    name: 'audit', 
    info: [
        { version: '1.52.100', path: 'swagger/audit/swagger1_52_100.yaml' },
        { version: '2.20.1', path: 'swagger/audit/swagger2_20_1.yaml' },
        { version: '3.0.0', path: 'swagger/audit/swagger3_0_0.yaml' }
    ]},{
    name: 'auth', 
    info: [
        { version: '0.0.1', path: 'swagger/auth/swagger0_0_1.yaml' },
        { version: '0.0.2', path: 'swagger/auth/swagger0_0_2.yaml' },
        { version: '1.0.0', path: 'swagger/auth/swagger1_0_0.yaml' }
    ]},{
    name: 'ess', 
    info: [
        { version: '2.0.0', path: 'swagger/ess/swagger2_0_0.yaml' }
    ]}
];


(async()=>{
    try{
        logger.debug(serviceTableName);
        logger.debug(versionTableName);

        const serviceParams = { RequestItems: { [serviceTableName]: [] } }
        serviceInfo.forEach(function (v, index) {
            serviceParams['RequestItems'][serviceTableName].push({
                PutRequest: {
                    Item: {
                        id: getUuid(v.name, 5),
                        servicename: v.name,
                        latestversion: '1.0.0',
                        lastupdated: (new Date()).getTime()
                    }
                }
            });
        })
        logger.debug(serviceParams);


        const versionParams = { RequestItems: { [versionTableName]: [] } }
        serviceInfo.forEach(function (v, index) {
            const id = getUuid(v.name, 5);
            v.info.forEach( ( ver )=>{
                versionParams['RequestItems'][versionTableName].push({
                    PutRequest: {
                        Item: {
                            id: id,
                            version: ver.version,
                            path: ver.path,
                            lastupdated: (new Date()).getTime()
                        }
                    }
                });
            })
        })
        
        logger.debug(versionParams);

        const serviceResult = await dynamodb.batchWrite(serviceParams).promise();
        logger.debug(serviceResult);

        const versionResult = await dynamodb.batchWrite(versionParams).promise();
        logger.debug(versionResult);

    }catch(error){
        logger.error(error);
    }
})()