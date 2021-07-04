// table row highlight
document.querySelectorAll('table tbody tr').forEach(e => e.addEventListener('click', function () {
    if (e.classList.contains('selected'))
        return;
    document.querySelectorAll('table tbody tr').forEach(e => e.classList.remove('selected'));
    e.classList.add('selected');
}));

// get selected row
function selectedRow() {
    var row = document.querySelector('tr.selected');
    return row;
}

function selected(title, content) {
    var selected_row = selectedRow();
    if (selected_row == null) {
        notify(title, content);
        return null;
    }
    return selected_row;
}

function delete_confirm() {
    if (!window.confirm("确定要删除吗?")) {
        return false;
    }
    return true;
}

// toast (notify)
const {
    Toast
} = bootstrap;

function initToast(title, body) {
    var htmlMarkup = `
  <div aria-atomic="true" aria-live="assertive" class="toast bg-primary text-white position-absolute end-0 top-0 m-3"
  style="z-index:5" role="alert" id="myAlert">
      <div class="toast-header">
            <strong class="me-auto">` + title + `</strong>
            <small></small>
            <button aria-label="Close" class="btn-close" 
                    data-bs-dismiss="toast" type="button">
            </button>
      </div>
      <div class="toast-body">` + body + ` </div>
  </div>
`;

    var template = document.createElement('template')
    html = htmlMarkup.trim()
    template.innerHTML = html
    return template.content.firstChild
}

function toast(title, body) {
    var toastEl = initToast(title, body);
    document.body.appendChild(toastEl)
    const myToast = new Toast(toastEl);
    myToast.show();
}

function action(add, edit, del) {
    $('#add').click(function () {
        window.location = add.url;
    });
    $('#edit').click(function () {
        var row = selected("提示", "请选择要更改的行");
        if (row != null) {
            var id = row.cells[edit.id].innerText;
            window.location = edit.url + id;
        }
    });
    $('#delete').click(function () {
        var row = selected("提示", "请选择要删除的行");
        if (row != null) {
            if (delete_confirm()) {
                var id = row.cells[del.id].innerText;
                ajax("DELETE", del.url + id, del.data, del.returnUrl);
            }
        }
    })
}

// method: PUT|DELETE
// url: endpoint
// data: the data send to server
// redirect: where to go after success
function ajax(method, url, data, redirect) {
    var xhr = new XMLHttpRequest();
    xhr.open(method, url, true);
    xhr.send(data);
    xhr.onload = function (e) {
        if ((xhr.status >= 200 && xhr.status < 300) || xhr.status == 304) {
            console.log(xhr.responseText);
            window.location = redirect;
        } else {
            console.log("What:", e)
        }
    };
    xhr.onerror = function (e) {
        console.log(e);
    };
}

// IP format validation
function isIP(strIP) {
    var re = /^(\d+)\.(\d+)\.(\d+)\.(\d+)|\:(\d*)$/g //匹配IP地址+端口的正则表达式  
    // TODO: 端口必须为数字，若是字母则RegExp.$5="" (10.19.1.8:abc)
    if (re.test(strIP)) {
        if (RegExp.$1 < 256 && RegExp.$2 < 256 && RegExp.$3 < 256 && RegExp.$4 < 256 && RegExp.$5 < 65535) return true;
    }
    return false;
}

function today() {
    var date = new Date();
    var month = date.getMonth() + 1;
    if (month < 10) {
        month = '0' + month.toString();
    }
    var day = date.getDate();
    if (day < 10) {
        day = '0' + day.toString();
    }
    return date.getFullYear() + "-" + month + "-" + day;
}

function notify(message, type) {
    var delay = 1000;
    if (type != "success") {
        delay = 60000;
    }
    return $.notify({
        message: message
    }, {
        type: type,
        placement: {
            from: "top",
            align: "center"
        },
        delay: delay,
        mouse_over: 'pause',
        animate: {
            enter: "animate__animated animate__fadeInDown",
            exit: "animate__animated animate__fadeOutUp"
        },
    });
}