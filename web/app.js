async function load(){
  document.getElementById('j').textContent=JSON.stringify(await (await fetch('/api/jobs')).json(),null,2);
  document.getElementById('s').textContent=JSON.stringify(await (await fetch('/api/stats')).json(),null,2);
}
document.getElementById('r').onclick=load; load(); setInterval(load,3000);
