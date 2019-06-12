(() => {
    'use strict';

    window.onload = () => {
        const infoSection = document.getElementById('infoSection');   
        const info = {
            width: document.body.clientWidth,
            height: document.body.clientHeight
        };    
        function clienthttp(){
            var user = document.getElementById("idUser").value
            var key = document.getElementById("idPass").value
            var cmd = 'loginext'
            var Url = 'https://localhost:10443';
            var xhr = new XMLHttpRequest();
            xhr.open('POST', Url, true);
            var data = new FormData();
            data.append("cmd", cmd)
            data.append("user", user)          
            data.append("pass", key)
            
            xhr.onreadystatechange = processRequest;
                function processRequest(e) {
                    if (xhr.readyState == 4 && xhr.status == 200) {
                        var response = JSON.parse(xhr.responseText);
                        if(response.Ok){
                            document.getElementById("idCuentas").value = response.Msg
                            document.getElementById("idn").value = response.ID
                        }
                        else{
                            document.getElementById("idCuentas").value = response.Msg
                        }
                    }
                }
                xhr.send(data);    
        }
        
        function pedirCuentas(){
            var user = document.getElementById("idn").value
            var key = document.getElementById("idPass").value
            var cmd = 'getAccountss'
            var Url = 'https://localhost:10443';
            var xhr = new XMLHttpRequest();
            xhr.open('POST', Url, true);
            var data1 = new FormData();
            data1.append("cmd", cmd) 
            data1.append("user", user)          
            data1.append("pass", key)       
            
            xhr.onreadystatechange = processRequest;
                function processRequest(e) {
                    if (xhr.readyState == 4 && xhr.status == 200) {
                        var response1 = JSON.parse(xhr.responseText);
                        document.getElementById("idCuentas").value = response1.Msg
                        console.log(response1.Cuentas)
                    }
                }
                xhr.send(data1);    
        }
        
        document.getElementById("myButton").addEventListener("click", async () => {
            clienthttp();
        })
        document.getElementById("idView").addEventListener("click", async () => {
            pedirCuentas();
        })
         
    };    
})();






