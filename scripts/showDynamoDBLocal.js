const AWS = require('aws-sdk');

const defaultIntervalTime = 10;
const defaultPort = 8000;

const program = require('commander');
program.version('0.0.1')
    .usage('[options] <file ...>')
    .option('-t, --time <n>', `interval time[sec](default: ${defaultIntervalTime})`, parseInt)
    .option('-p, --port <n>', `port(default: ${defaultPort})`, parseInt)
    .parse(process.argv);

const sleep = sec => new Promise(resolve => setTimeout(resolve, sec*1000));
const intervalTime = program.time || defaultIntervalTime;
const port = program.port || defaultPort;
const endpoint = `http://localhost:${port}`;

AWS.config.update({ region: 'sa-east-1' });
const dynamodb = new AWS.DynamoDB({ endpoint });
const documentClient = new AWS.DynamoDB.DocumentClient({ endpoint });


const showTable = async (tableList) => {
    for (let i = 0; i < tableList.length; i = i + 1) {
        console.log(`\n++++++++++++++++++++ [${i}] ${tableList[i]}++++++++++++++++++++`);
        const table = await documentClient.scan({ TableName: tableList[i] }).promise();
        console.table(table.Items);
    }
}

(async()=>{
    try {

        const tableList = (await dynamodb.listTables({}).promise()).TableNames;


        if('time' in program) {
            console.clear();
            while(true){
                console.log(`Endpoint: ${endpoint}`);
                await showTable(tableList);
                await sleep(intervalTime);
                console.clear();
            }
        } else {
            console.log(`Endpoint: ${endpoint}`);
            await showTable(tableList);
        }
        
    }catch(error){
        console.error(error);
    }
})()