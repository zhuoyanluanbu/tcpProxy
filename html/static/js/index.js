//配置数据
var testData = [
    // {id: 1, port1: 1, ip: '123.133.434', port2: 123, tls:false, crtPath:'xxx', keyPath:'kkk'},
    // {id: 2, port1: 2, ip: '123.143.424', port2: 123343, tls:false, crtPath:'yyy', keyPath: 'lll'}
];

//当前连接
var testData2 = [
    // {id: 1, port1: 1, ip1: '123.133.434', port2: 123, ip2: '123.32.23.4'},
    // {id: 2, port1: 2, ip1: '123.143.424', port2: 123343, ip2: '123.43.43.4'}
];

//文件名
var file_test = ['1233', '1231313'];

var size = 0;

var removeArr = [];

//----------------------------------------------------------------------------------------
//-------------------------------------加载配置---------------------------------------------------
function loadConfigs() {
    /*增删改查页面加载*/
    var html = '';
    for (var i = 0; i < testData.length; i++) {
        var td = testData[i]
        var checked = td.tls?"checked":""
        html +=
            '<tr>' +
                '<td>' +
                    '<input onblur="updateCheck(' + i + ',this)" type="checkbox" class="encode-checkbox"'+ checked +'></td>' +
                '<td>' +
                    '<input class="input-port" onchange="updatePort1(' + i + ',this)" id="port1-' + i + '" style="width: 3.8em;text-align: center" type="number" value="' + td.port1 + '">' +
                '</td>' +
                '<td style="text-align: center;">' +
                    '<input style="width: 12.6em" class="input-ip" id="ip-' + i + '" onblur="updateIp(' + i + ',this)" type="text" value="' + td.ip + '"/>' +
                    // ' : ' +
                    // '<input id="port2-' + i + '" onblur="updatePort2(' + i + ',this)" style="width: 3.8em" type="number" value="' + td.port2 + '"/>' +
                '</td>' +
                '<td>' +
                    // '<div style="font-size: 0.4em;">' +
                    //     '<div style="line-height: 1.2em">' + td.crtPath.replace('./certs/','') + '</div>' +
                    //     '<div style="line-height: 1.2em">' + td.keyPath.replace('./certs/','') + '</div>' +
                    // '</div>' +
                    '<button class="manageBtn" type="button" onmouseleave="hideBindedCerts()" onmouseover="showBindedCerts('+i+')" onclick="chooseFile(' + i + ')">选择证书</button>' +
                '</td>' +
                '<td>' +
                    '<button class="manageBtn" style="background-color: #d35400" type="button" onclick="deleteData(this,' + i + ')">删除</button>' +
                '</td>' +
            '</tr>';
    }
    $('#leftTbody').html(html);
    size = testData.length;
};

function mouseXY(ev){
    if(ev.pageX || ev.pageY){
        return {x:ev.pageX, y:ev.pageY};
    }
    return {
        x:ev.clientX + document.body.scrollLeft - document.body.clientLeft,
        y:ev.clientY + document.body.scrollTop - document.body.clientTop
    };
}
function showBindedCerts(i) {
    var td = testData[i]
    if (td.crtPath == '' || td.keyPath == ''){
        return
    }
    var crt = td.crtPath.replace('./certs/','')
    var key = td.keyPath.replace('./certs/','')

    var xy = mouseXY(window.event)
    var top = xy.y+5 + "px"
    var left = xy.x+5 + "px"
    // alert(top+"---"+left)
    var div =
        "<div " +
            "id='certsAlert' " +
            "style='" +
                "background-color: #dcddff;" +
                "border: 1px solid #a299ff;" +
                "font-size: 0.7em;" +
                "padding: 0.5em 0.8em 0.7em 0.8em;" +
                "color: black;" +
                "border-radius: 0.3em;" +
                "width: 10em;" +
                "position: absolute;" +
                "z-index: 9999;" +
                "top:"+top+"; " +
                "left:"+left+"'>" +
            "证书："+crt +
            "<br/>" +
            "私钥：" + key
        "</div>"
    $("body").append(div)
}
function hideBindedCerts() {
    $('#certsAlert').remove()
}

function loadConfigsFromApi() {
    $.ajax({
        url: '/getConfig',
        success: function(res){
            testData = [];
            res = JSON.parse(res);
            for (var f in res){
                var data = res[f];
                var source = data.source;
                var destination = data.destinations;
                var port1 = source.split(":")[0];
                var ip = destination;
                var tls = data.tls;
                var crtPath = '';
                var keyPath = '';
                if (tls) {
                    crtPath = data.tlsCf.crtPath;
                    keyPath = data.tlsCf.keyPath;
                }
                var obj = {
                    port1:port1,
                    ip:ip,
                    // port2: port2,
                    tls: tls,
                    crtPath: crtPath,
                    keyPath: keyPath
                };
                testData.push(obj)
            }
            loadConfigs()
        }
    })
}
loadConfigsFromApi();

//----------------------------------------------------------------------------------------
//-------------------------------------加载当前连接信息---------------------------------------------------
function loadConns() {
    var html1 = '';
    for (var i = 0; i < testData2.length; i++) {
        html1 +=
            '<tr>' +
            '<td>' +
            '<span>' + testData2[i].ip1 + '</span>' + ' : <span>' + testData2[i].port1 + '</span>' +
            '</td>' +
            '<td>' +
            '<span> → </span>' +
            '</td>' +
            '<td>' +
            '<span>' + testData2[i].ip2 + '</span>' +
            ' : ' +
            '<span>' + testData2[i].port2 + '</span>' +
            '</td>' +
            '<td>' +
            '<button class="manageBtn" style="background-color: #d35400" type="button" onclick="drop(this,\'' + testData2[i].id + '\')">踢掉</button>' +
            '</td>' +
            '</tr>';
    }
    $('#rightTbody').html(html1);
}

function loadConnsFromApi() {
    $.ajax({
        url: '/getAliveConns',
        success: function(res){
            testData2 = []
            res = JSON.parse(res)
            //[{"id":"xxx","remote":"xxxx.xxx.xxx.xxx:yyyy","local":"xxxx.xxx.xxx.xxx:yyyy"}]
            for (var d in res){
                var data = res[d]
                var id = data.id
                var remote = data.remote
                var local = data.local
                var port1 = remote.split(":")[1]
                var ip1 = remote.split(":")[0]
                var port2 = local.split(":")[1]
                var ip2 = local.split(":")[0]
                var obj = {
                    id:id,
                    port1:port1,
                    ip1:ip1,
                    port2:port2,
                    ip2: ip2
                }
                testData2.push(obj)
            }
            console.log(testData2)
            loadConns()
            setTimeout(loadConnsFromApi,2000)
        },
        error: function(err){
            setTimeout(loadConnsFromApi,5000)
        }
    })

}
loadConnsFromApi()

//----------------------------------------------------------------------------------------
//------------------------------/*绑定文件页面加载*/----------------------------------------------------------
function loadFilesNames() {
    /*绑定文件页面加载*/
    var fileHtml = '';
    for (var i in file_test) {
        fileHtml +=
            '<tr>' +
                '<td>' +
                    '<input type="checkbox" name="file-checkbox" value="' + file_test[i] + '"  />' +
                    '<span style="margin-left: 2em">' + file_test[i] + '</span>' +
                '</td>' +
            '</tr>';
    }
    $('#choose-file-table').html(fileHtml);
}

function loadFilesNamesFromApi() {
    $.ajax({
        url: '/getCertFileNames',
        success: function(res){
            file_test = []
            res = JSON.parse(res)
            for (var f in res){
                file_test.push(res[f])
            }
            setTimeout(loadFilesNamesFromApi,120000)
        },
        error: function (err) {
            setTimeout(loadFilesNamesFromApi,300000)
        }
    })
}
loadFilesNamesFromApi()
//------------------------------------------------------------------
//------------------------------------------------------------------



/*删除缓存中的文件*/
function deleteData(btn, sort) {
    if (!flag) {
        $(btn).parent().parent().remove();
        removeArr.push(sort);
    }else alert("程序运行中，请先停止！");
}

/*调用接口删除文件*/
function drop(btn, id) {
    var isOk = confirm("确定踢掉这个链接?");
    if (!isOk){
        return
    }
    $.ajax({
        url:'/dropConn?id='+id,
        success: function (res) {
            if (flag){
                $(btn).parent().parent().remove();
            }
        }
    })
}


/*增加数据*/
function addData() {
    if (flag) {
        alert("程序运行中，请先停止！");
        return
    }
    if (!checkInput()) return;
    var html =
        '<tr>' +
        '<td>' +
        '<input type="checkbox" class="encode-checkbox">' +
        '</td>' +
        '<td>' +
        '<input class="input-port" onchange="updatePort1(' + size + ',this)" id="port1-' + size + '" style="width: 3.8em;text-align:center" type="number" />' +
        '</td>' +
        '<td style="text-align: center">' +
        '<input style="width:8em;" id="ip-' + size + '" onblur="updateIp(' + size + ',this)" type="text" class="input-ip" />' +
        // ' : ' +
        // '<input id="port2-' + size + '" onblur="updatePort2(' + size + ',this)" style="width: 3.8em" type="number" />' +
        '</td>' +
        '<td>' +
        '<button class="manageBtn" type="button" onclick="chooseFile(' + size + ')">选择证书</button><' +
        '/td>' +
        '<td>' +
        '<button class="manageBtn" style="background-color: #d35400" type="button" onclick="deleteData(this,' + size + ')">删除</button>' +
        '</td>' +
        '</tr>';
    $('#leftTbody').append(html);
    size++;
    if (!checkInput()) return;
    var obj = {
        port1:0,
        ip:"",
        tls: false,
        crtPath: "",
        keyPath: ""
    };
    testData.push(obj)
}

function showUpload() {
    $('.topDiv').css('display', 'flex');
    $('#upload-alert').css('display', 'flex');
    $('#choose-alert').css('display', 'none');
}

/*上传文件*/
function uploadFile() {
    var files1 = $('#file1')[0].files; //获取上传的文件列表
    var files2 = $('#file2')[0].files; //获取上传的文件列表

    if (files1.length > 0){
        var formData1 = new FormData(); //新建一个formData对象
        formData1.append("file", files1[0]); //append()方法添加字段
        upload(formData1,files1[0].name);
    }

    if (files2.length > 0){
        var formData2 = new FormData(); //新建一个formData对象
        formData2.append("file", files2[0]); //append()方法添加字段
        upload(formData2,files2[0].name);
    }
}

function upload(formData,filename) {
    $.ajax({
        url: '/test',
        type:  'POST',
        data: formData,
        processData: false,
        headers:{'Content-Type':'application/json;charset=utf8','filename':filename},
        success: function (data) {

        }
    })
}

function hideAlert() {
    $('.topDiv').css('display', 'none');
}

var flag = false;

function getRunState() {
    $.ajax({
        url: '/runState',
        success: function (res){
            if (res == 'true') {
                flag = true
                $('#applyBtn').hide();
                $('#stopBtn').show();

                $('.input-port').attr('readOnly','readOnly');
                $(".input-ip").attr('readOnly','readOnly');
                $(".encode-checkbox").attr('disabled','disabled');
            }else {
                flag = false;
                $('#applyBtn').show();
                $('#stopBtn').hide();

                $('.input-port').removeAttr('readOnly');
                $(".input-ip").removeAttr('readOnly');
                $(".encode-checkbox").removeAttr('disabled');
            }
        },
        error: function (res) {
        }
    })
}
getRunState();

function startRun() {
    if (!flag) {
        if (removeArr.length > 0) {
            for (var i = 0; i < removeArr.length; i++) {
                testData.splice(removeArr[i], 1);
            }
        }
        if (!checkInput()) return
        console.log(JSON.stringify(testData));
        $.ajax({
            url: '/startup',
            type: 'post',
            data: JSON.stringify(testData),
            success: function (res){
                if (res == 'ok'){
                    flag = true;
                    $('#applyBtn').hide();
                    $('#stopBtn').show();
                    alert("...启动成功...");
                    //启动之后禁用
                    $('.input-port').attr('readOnly','readOnly');
                    $(".input-ip").attr('readOnly','readOnly');
                    $(".encode-checkbox").attr('disabled','disabled');
                }else {
                    alert(res);
                }

            },
            error: function (err) {
                alert(JSON.stringify(err))
            }
        })
    }
}

function checkInput() {
    for (var i=0;i<testData.length;i++){
        var t = testData[i];
        if (t.port1 <= 0){
            alert("请检查第["+(i+1)+"]排的监听端口")
            return false
        }
        t.ip = t.ip.replace(/，/g,',').replace(/;/g,',').replace(/；/g,',')
        // if (t.port2 <= 0){
        //     alert("请检查第["+(i+1)+"]排的目的地址端口")
        //     return false
        // }

        // if ( t.ip != 'localhost' && (t.ip == "" || t.ip.split(".").length!=4)){
        //     alert("请检查第["+(i+1)+"]排的目的地址ip")
        //     return false
        // }

        // if (t.port1 == t.port2 && (t.ip == 'localhost' || t.ip == '127.0.0.1')) {
        //     alert("请检查第["+(i+1)+"]排的数据，不能和本机的监听端口相同")
        //     return false
        // }
    }
    return true
}

function stopRun() {
    if (flag){
        $.ajax({
            url: '/shutdown',
            success: function (res){
                flag = false;
                $('#applyBtn').show();
                $('#stopBtn').hide();
                //停止之后启用input
                $('.input-port').removeAttr('readOnly');
                $(".input-ip").removeAttr('readOnly');
                $(".encode-checkbox").removeAttr('disabled');
            },
            error: function (res) {
            }
        })
    }
}

function updateCheck(sort, text) {
    if (!flag) {
        testData[sort].tls = $(text).prop('checked');
    }
}

function updatePort1(sort, text) {
    for (var t in testData){
        // testData[t].port1 = testData[t].port1.replace(/，/g,',');
        if ($(text).val() == testData[t].port1) {
            alert("监听端口重复")
            $(text).val('');
            return
        }
    }

    if (!flag) {
        testData[sort].port1 = $(text).val();
    }
}

function updateIp(sort, text) {
    if (!flag) {
        testData[sort].ip = $(text).val();
    }
}

function updatePort2(sort, text) {
    if (!flag) {
        testData[sort].port2 = $(text).val();
    }
}


function chooseFile(sort) {
    loadFilesNames()
    if (!flag) {
        $('.topDiv').css('display', 'flex');
        $('#choose-alert').css('display', 'flex');
        $('#upload-alert').css('display', 'none');

        $('#dataSort').val(sort);
    } else {
        alert("程序运行中，请先停止！");
        return;
    }
}

/*绑定文件*/
function bindFile() {
    var sort = $('#dataSort').val();
    var files = {crtPath:"",keyPath:""};

    var checkedFiles = $('input[name="file-checkbox"]:checked')
    if (checkedFiles.length > 2){
        alert("只能选择一个key，一个pem")
        return
    }
    checkedFiles.each(function () {
        var name = $(this).val()
        if (name.indexOf(".") > 0){
            var sux = name.split(".")[1]
            if (sux == "key"){
                files.keyPath = name
            }else if (sux == "pem") {
                files.crtPath = name
            }
        }
    });
    if (files.keyPath == ""){
        alert("请选择key")
        return
    }
    if (files.crtPath == ""){
        alert("请选择pem")
        return
    }
    testData[sort].keyPath = files.keyPath;
    testData[sort].crtPath = files.crtPath;
    hideAlert();
}