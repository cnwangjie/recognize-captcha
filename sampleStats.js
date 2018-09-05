const fs = require('fs')
const path = require('path')
const files = fs.readdirSync('/home/wangjie/Workspace/tmp/zfc-samples')
const sum = files.length

const letterAnalytics = {}

files.map(file => {
    const code = path.basename(file, '.png')
    code.split('').map(letter => {
        if (letterAnalytics[letter]) {
            letterAnalytics[letter] += 1
        } else {
            letterAnalytics[letter] = 1
        }
    })
})

console.log(`file sum: ${sum}\n`)

const letters = '0123456789abcdefghijklmnopqrstuvwxyz'.split('').filter(i => {
    // return letterAnalytics[i]
    return true
}).map(i => {
    return `${i}:${letterAnalytics[i] ? letterAnalytics[i] : '0 <---'}`
}).join('\n')
console.log(`${letters}`)
