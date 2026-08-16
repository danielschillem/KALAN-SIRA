const TOKEN_KEY='kalan_sira_session';
export const session={get:()=>{try{return JSON.parse(localStorage.getItem(TOKEN_KEY)||'null')}catch{return null}},set:v=>localStorage.setItem(TOKEN_KEY,JSON.stringify(v)),clear:()=>localStorage.removeItem(TOKEN_KEY)};
async function request(path,options={}){const s=session.get();const headers={'Content-Type':'application/json',...(options.headers||{})};if(s?.token)headers.Authorization=`Bearer ${s.token}`;const r=await fetch(path,{...options,headers});const data=await r.json().catch(()=>({}));if(!r.ok)throw new Error(data?.error?.message||'Erreur serveur');return data}
export const api={
 requestOTP:phone=>request('/api/v1/auth/otp/request',{method:'POST',body:JSON.stringify({phone})}),
 verifyOTP:(challenge_id,code)=>request('/api/v1/auth/otp/verify',{method:'POST',body:JSON.stringify({challenge_id,code})}),
 dashboard:()=>request('/api/v1/school/dashboard'),
 admissionCatalog:()=>request('/api/v1/admissions/catalog'),
 admissionFees:classID=>request(`/api/v1/admissions/fees/${encodeURIComponent(classID)}`),
 createAdmission:payload=>request('/api/v1/admissions',{method:'POST',body:JSON.stringify(payload)})
};
