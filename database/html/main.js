var path = window.location.pathname;
if (path=="/static/show.html") {
  path = "/show";
}
fetch(path, {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then(response => {
  if (!response.ok) {
    throw new Error('Network response was not ok');
  }
  return response.json();
})
.then(data => {
  var memos=data
  const memoList = document.getElementById('memoList');

memos.forEach(memo => {
  const li = document.createElement('li');
  li.className = 'memo-item';
  li.textContent = memo.Body;
  li.setAttribute('data-uid', memo.Uid);
  memoList.appendChild(li);

  let timer;
  li.addEventListener('mousedown', () => {
    timer = setTimeout(() => {
      const action = confirm("选择操作: 确认删除, 取消更新");
      if (action) {
        deleteMemo(memo.Uid); // 删除操作
      } else {
        const newBody = prompt("请输入新的备忘内容", memo.Body);
        if (newBody) updateMemo(memo.Uid, newBody); // 更新操作
      }
    }, 1000); // 长按时间阈值，例如1000毫秒
  });

  li.addEventListener('mouseup', () => {
    clearTimeout(timer); // 取消长按事件
  });
});

function deleteMemo(uid) {
  console.log(`删除Uid ${uid}`);
  
  // 实际项目中这里应发送POST请求到服务器端点进行删除
}

function updateMemo(uid, newBody) {
  console.log(`更新Uid ${uid}，新内容：${newBody}`);
  // 实际项目中这里应发送POST请求到服务器端点进行更新
}
})
.catch(error => {
  // 捕获任何网络请求错误
  console.error('There was a problem with the fetch operation: ', error);
});