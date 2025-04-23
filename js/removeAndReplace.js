// 移除不需要的元素 以及 替换复制不了的元素
// 移除页头
function removeHead() {
    document.querySelector('#mw-head').remove();
}

// 移除页尾
function removeFooter() {
    document.querySelector('#cpp-footer-base').remove();
}

// 移除 H1
function removeH1() {
    document.querySelector('#firstHeading').remove();
}

// 移除 面包屑导航菜单
function removeBreadcrumbMenu() {
    document.querySelector('#mw-content-text div.t-navbar').remove();
}

// 替换定义所在表格
// 注意必须在其他替换之前执行
function replaceDefineTableToPreCode() {
    // 获取第一个 table.t-dcl-begin 元素
    const table = document.querySelector('table.t-dcl-begin');

    if (table) {
        // 用于存储转换后的内容
        let preContent = '';

        // 遍历表格的行，除了表头行
        for (let i = 0; i < table.rows.length; i++) {
            const row = table.rows[i];
            const tdContent = row.cells[0].textContent.trim();
            const comment1 = row.cells[1].textContent.trim();
            const comment2 = row.cells[2].textContent.trim();
            let line
            if (tdContent.startsWith("在标头")) {
                // line = "<span style='display:flex;'><span>//" + tdContent + "<span><span>";
                line = `<span style='display:flex;'><span>// ${tdContent}<span><span>`;
            } else {
                line = `<span style='display:flex;'><span>${tdContent}<span><span>`;
                // line = "<span style='display:flex;'><span>" + tdContent + "<span><span>";
                if (comment1) {
                    if (comment2) {
                        line += "// " + comment1 + comment2 + "<span></span>";
                    } else {
                        line += "// " + comment1 + "<span></span>";
                    }
                } else {
                    if (comment2) {
                        line += "// " + comment2 + "<span></span>";
                    }
                }
            }
            line += '</span>\n';
            preContent += line;
        }

        // 用 <pre> 标签包裹内容
        const result = "<pre>" + preContent + "</pre>";

        // 创建一个新的元素来存储处理后的内容
        const newElement = document.createElement('div');
        newElement.innerHTML = result;

        // 移除原表格
        table.parentNode.replaceChild(newElement, table);
    }
}


// 移除 h3=引用 的内容
function removeH3YinYong() {
    const allH3 = document.querySelectorAll('h3');

// 遍历查找包含"引用"的 h3 标签
    let targetH3 = null;
    for (const h3 of allH3) {
        const headlineSpan = h3.querySelector('.mw-headline');
        if (headlineSpan && headlineSpan.textContent.trim() === '引用') {
            targetH3 = h3;
            break;
        }
    }

    if (targetH3) {
        // 删除操作
        let nextElement = targetH3.nextElementSibling;

        // 循环处理后续元素，直到遇到下一个 h3
        while (nextElement && nextElement.tagName !== 'H3') {
            const currentElement = nextElement;
            nextElement = currentElement.nextElementSibling; // 先保存下一个元素

            // 删除符合要求的 div
            if (currentElement.tagName === 'DIV' && currentElement.className.startsWith('t-ref-std')) {
                currentElement.remove();
            }
        }

        // 最后删除目标 h3 本身
        targetH3.remove();
    }
}


// 移除 h3=参阅 的内容
function removeH3CanYue() {
    const allH3 = document.querySelectorAll('h3');

// 遍历查找包含"引用"的 h3 标签
    let targetH3 = null;
    for (const h3 of allH3) {
        const headlineSpan = h3.querySelector('.mw-headline');
        if (headlineSpan && headlineSpan.textContent.trim() === '参阅') {
            targetH3 = h3;
            break;
        }
    }

    if (targetH3) {
        // 删除操作
        let nextElement = targetH3.nextElementSibling;

        // 循环处理后续元素，直到遇到下一个 h3
        while (nextElement) {
            const currentElement = nextElement;
            nextElement = currentElement.nextElementSibling; // 先保存下一个元素

            // 删除
            currentElement.remove();
        }

        // 最后删除目标 h3 本身
        targetH3.remove();
    }
}

function removeRunCode() {
    const allDiv = document.querySelectorAll('.t-example-live-link');
    for (const div of allDiv) {
        div.remove()
    }
}

// 移除页面中所有 <p><br></p>
function removePBr() {
    const paragraphs = document.getElementsByTagName('p');
    const paragraphsToRemove = [];

    for (let i = 0; i < paragraphs.length; i++) {
        const paragraph = paragraphs[i];
        if (paragraph.childElementCount === 1 && paragraph.firstElementChild.tagName === 'BR') {
            paragraphsToRemove.push(paragraph);
        }
    }

    paragraphsToRemove.forEach(paragraph => {
        paragraph.parentNode.removeChild(paragraph);
    });
}

function replaceSpanSourceC() {
    document.querySelectorAll('span.mw-geshi.c.source-c').forEach(span => {
        if (!span.querySelector(".kw486")) {
            // 在原有内容前后添加反引号（保留内部HTML结构）
            span.innerHTML = "\u0060" + span.innerHTML + "\u0060";
        }
    });
}

function replaceDlToUl() {
    document.querySelectorAll('dl').forEach(dl => {
        const ul = document.createElement('ul');
        ul.style.listStyleType = 'disc'; // 保持列表样式

        // 将 <dd><ul> 转换为常规列表
        dl.querySelectorAll('li').forEach(li => {
            const newLi = document.createElement('li');
            newLi.innerHTML = "&zeroWidthSpace;" + li.innerHTML; // 保留原有内容
            ul.appendChild(newLi);
        });

        // 替换原有 <dl> 结构
        dl.parentNode.replaceChild(ul, dl);
    });
}

function replaceXuHao() {
    // 找到所有 class="t-li" 的元素
    const elements = document.querySelectorAll('.t-li1 .t-li');

    // 遍历每个元素
    elements.forEach(element => {
        // 获取元素的文本内容
        const text = element.textContent.trim();

        // 使用正则表达式检查内容是否符合要求
        const regex = /^\d+\)$/; // 匹配以数字开头，以 ) 结尾的内容

        if (regex.test(text)) {
            // 替换 ) 为 ）
            element.textContent = text.replace(')', '）');
        }
    });
}

// 在p标签前面添加一个span，其内容为&zeroWidthSpace;用于后期在markdown中替换成Tab符号
function replaceP() {
    // 在p标签前面插入&zeroWidthSpace;
    let pElements = document.querySelectorAll('p');

    pElements.forEach(function (element) {
        let newSpan = document.createElement('span');
        newSpan.textContent = '&zeroWidthSpace;';
        if (element.firstChild) {
            element.insertBefore(newSpan, element.firstChild);
        } else {
            // 如果 p 元素没有子节点，直接将新 span 元素添加到 p 元素中
            element.appendChild(newSpan);
        }
    });
}

function replaceTableContentAddBr() {
    // 增加 .t-lines span 下的文本内容，增加@!br !@用于后面在md文档中替换成换行符
    let spans = document.querySelectorAll('.t-lines span');

// 遍历每个 span 元素
    spans.forEach(span => {
        // 获取原文本内容
        const originalText = span.textContent;
        // 在原文本末尾添加两个空格
        span.textContent = originalText + '@!br /!@';
    });
}

function replaceLtGt() {
    // 替换tt标签中的< 和 > 为 @!和 !@
    let ttElements = document.querySelectorAll('a tt');
    ttElements.forEach((tt) => {
        let text = tt.textContent;
        if (text.length > 0) {
            text = '@!' + text.slice(1, -1) + '!@';
            tt.textContent = text;
        }
    });
}

function replaceNotice() {
    const tables = document.querySelectorAll('table.metadata.plainlinks.ambox.mbox-small-left.ambox-notice');

    tables.forEach(table => {
        const blockquote = document.createElement('blockquote');
        const rows = table.querySelectorAll('tr');

        rows.forEach(row => {
            const p = document.createElement('p');
            let hasContent = false;
            const cells = Array.from(row.querySelectorAll('td'));

            for (let i = 0; i < cells.length; i++) {
                const cell = cells[i];
                if (cell.textContent.trim() !== '') {
                    hasContent = true;
                    const clonedCell = cell.cloneNode(true);
                    const brTags = clonedCell.querySelectorAll('br');
                    brTags.forEach(br => {
                        const newBr = document.createElement('br');
                        br.parentNode.insertBefore(newBr, br.nextSibling);
                    });
                    clonedCell.childNodes.forEach(child => {
                        p.appendChild(child.cloneNode(true));
                    });
                }
            }

            if (hasContent) {
                blockquote.appendChild(p);
            }
        });

        table.parentNode.replaceChild(blockquote, table);
    });
}

// 替换形似行的table，将table的内容替换成p标签包裹
function replaceTRevBeginTable() {
    const tables = document.querySelectorAll('table.t-rev-begin');

    tables.forEach(table => {
        const parent = table.parentNode;
        const rows = table.querySelectorAll('tr');

        rows.forEach(row => {
            const p = document.createElement('p');
            const cells = row.querySelectorAll('td');

            cells.forEach(cell => {
                const clonedCell = cell.cloneNode(true);
                const brTags = clonedCell.querySelectorAll('br');
                brTags.forEach(br => {
                    br.parentNode.removeChild(br);
                });

                const pTags = clonedCell.querySelectorAll('p');
                if (pTags.length > 0) {
                    pTags.forEach(pTag => {
                        pTag.childNodes.forEach(child => {
                            p.appendChild(child.cloneNode(true));
                        });
                    });
                } else {
                    clonedCell.childNodes.forEach(child => {
                        p.appendChild(child.cloneNode(true));
                    });
                }

                // if (cells[cells.length - 1]!== cell) {
                //     p.appendChild(document.createElement('br'));
                // }
            });

            parent.insertBefore(p, table);
        });

        parent.removeChild(table);
    });
}

function replaceTableAddVersionAfterIdentifier() {
    // 替换 html中的td行内容: 即将C99版本加在函数名后面，避免出现多行C99在单独行的情况
    let divElements = document.querySelectorAll('div.t-dsc-member-div');
    divElements.forEach((div) => {
        const firstDiv = div.querySelector('div:first-child');
        const secondDiv = div.querySelector('div:nth-child(2)');

        if (firstDiv && secondDiv) {
            const secondDivSpans = secondDiv.querySelectorAll('.t-lines span');
            const firstDivSpans = firstDiv.querySelectorAll('.t-lines span');

            firstDivSpans.forEach((span, index) => {
                if (secondDivSpans[index]) {
                    let contentToAppend = secondDivSpans[index].textContent;
                    // console.log("span.textContent=",span.textContent,"contentToAppend=",contentToAppend)
                    span.textContent += "   " + contentToAppend;
                }
            });
            // 移除不要的html标签
            secondDiv.remove();
        }
    });
}

removeHead();
removeFooter();
removeH1();
removeBreadcrumbMenu();
removePBr();
replaceLtGt();
replaceDefineTableToPreCode();
// removeH3YinYong();
// removeH3CanYue();
removeRunCode();
replaceSpanSourceC();
replaceDlToUl();
replaceXuHao();
replaceP();
replaceTableContentAddBr();
replaceNotice();
replaceTRevBeginTable();
replaceTableAddVersionAfterIdentifier();

// 转义字符为正则表达式中可以直接使用的字符
function quoteMeta(str) {
    return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

function getAllHeaderInfo() {
    const pattern = /<([^<>./]+)\.h>/;
    const table = document.querySelector('table.t-dsc-begin');
    const linkData = [];
    // 用于记录已经存在的 header
    const existingHeaders = {};
    if (table) {
        // 获取表格的所有行
        const rows = table.querySelectorAll('tr.t-dsc');
        rows.forEach(row => {
            // 获取第一列的所有链接
            const firstColumn = row.querySelector('td:first-child');
            const links = firstColumn.querySelectorAll('a');
            links.forEach(link => {
                const fullHeader = link.textContent.trim();
                const headerMatch  = fullHeader.match(pattern);
                const url = link.href;
                const header = headerMatch ? headerMatch[1] : "";
                if (!existingHeaders[header]) {
                    linkData.push({
                        header: header,
                        fullHeader: fullHeader,
                        url: url
                    });
                    // 标记该 header 已经存在
                    existingHeaders[header] = true;
                }
            });
        });
        // 打印结果
        console.log(linkData);
    }
    return linkData
}

// 根据表格第二列中提示，返回类型标识
function calTypeFlag(str) {
    // (函数) -> f <-函数
    // (常量) (枚举) -> e <- 枚举
    // (宏常量) (关键词宏) (宏函数) -> m <- 宏
    // (typedef) (结构体) -> t <- 类型
    // 没有的情况 -> 暂时归类为 t <- 类型
    if(str === "(函数)") {
        return "f"
    }

    if(str === "(常量)" || str === "(枚举)") {
        return "e"
    }

    if(str === "(宏常量)" || str === "(关键词宏)"|| str === "(宏函数)") {
        return "m"
    }

    if(str === "(typedef)" || str === "(结构体)") {
        return "t"
    }

    return "t"
}


function GetIdentifierFromGaiYao() {
    let ids = []
    let exist = {}

    const divMwGeshis = document.querySelectorAll("div.mw-geshi")
    divMwGeshis.forEach(mw => {
        let typeNames = [];
        let macroNames = [];
        let functionNames = [];
        // 判断其紧邻的元素中的上一个p标签是否有如下内容作为开头：仅当实现定义了
        let remark = ""
        const prevE = mw.previousElementSibling;
        console.log("prevE.textContent=", prevE.textContent)
        if (prevE.tagName === "P" && prevE.textContent.startsWith("仅当实现定义了")) {
            remark = prevE.textContent
        }

        const pre = mw.querySelector("div.c.source-c > pre.de1")
        const preCode = pre.textContent
        const lines = preCode.split('\n');
        lines.forEach(line => {
            // 去除注释
            line = line.replace(/\/\/.*|\/\*[\s\S]*?\*\//g, '');
            line = line.trim();

            // 匹配类型名称
            const typeMatch = line.match(/^typedef\s+(?:struct|union|enum)?\s*(\w+);/);
            if (typeMatch) {
                typeNames.push(typeMatch[1]);
            }

            // 匹配宏名称
            const macroMatch = line.match(/^#define\s+(\w+)\b/);
            if (macroMatch) {
                macroNames.push(macroMatch[1]);
            }

            // 匹配函数名
            const functionMatch = line.match(/^([\w\s*_]+)\s+(\w+)\s*\(/);
            if (functionMatch) {
                // 排除类型定义、宏定义和关键字开头的情况
                if (!/^(typedef|#define|struct|union|enum)/.test(line)) {
                    functionNames.push(functionMatch[2]);
                }
            }
        });
        // console.log(1,typeNames)
        // console.log(1,macroNames)
        // console.log(1,functionNames)
        // console.log("---------------------------")
        typeNames = [...new Set(typeNames)]
        macroNames = [...new Set(macroNames)]
        functionNames = [...new Set(functionNames)]
        // console.log(2,typeNames)
        // console.log(2,macroNames)
        // console.log(2,functionNames)
        if (typeNames.length > 0) {
            typeNames.forEach(id => {
                exist[id] = true
                ids.push({
                    id: id,
                    typ: "t",
                    url: "",
                    remark: remark,
                    dsc: "",
                })
            })
        }

        if (macroNames.length > 0) {
            macroNames.forEach(id => {
                exist[id] = true
                ids.push({
                    id: id,
                    typ: "m",
                    url: "",
                    remark: remark,
                    dsc: "",
                })
            })
        }

        if (functionNames.length > 0) {
            functionNames.forEach(id => {
                exist[id] = true
                ids.push({
                    id: id,
                    typ: "f",
                    url: "",
                    remark: remark,
                    dsc: "",
                })
            })
        }

        const aLinks = mw.querySelectorAll("div.c.source-c > pre.de1 a")
        aLinks.forEach(a => {
            for(const obj of ids) {
                if(obj.id === a.textContent.trim()) {
                    obj.url = a.href
                }
            }
        })
    })

    console.log(ids)

    const trs = document.querySelectorAll("div.t-dsc-begin tr.t-dsc")
    trs.forEach(tr => {
        const firstTd = tr.querySelector("td:first-child")
        const aE = firstTd.querySelector("a")
        if (aE) {
            const spans = firstTd.querySelector("div:first-child span.t-lines span")
            const typFlagSpan = tr.querySelector("td:nth-child(2) span.t-mark")
            spans.forEach(span => {
                const id = span.textContent.trim()
                if(!exist[id]) {
                    ids.push({
                        id: id,
                        typ: calTypeFlag(typFlagSpan.textContent.trim()),
                        url: aE.href,
                        remark: "",
                        dsc: "",
                    })
                    exist[id] = true
                } else {
                    for(const obj of ids) {
                        if (obj.id === id) {
                            obj.url = aE.href
                            obj.dsc = ""
                        }
                    }
                }
            })
        } else {
            const codeB = firstTd.querySelector("code > b")



        }


    })


}



