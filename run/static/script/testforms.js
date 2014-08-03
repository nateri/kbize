/*
 $Id: testforms.js,v 1.6 2008/07/04 10:22:56 agoubard Exp $

 Copyright 2003-2008 Online Breedband B.V.
 See the COPYRIGHT file for redistribution and use restrictions.
*/

function doHandleSubmit(form)
{
   var elemInSearch = document.getElementById('in_search');
   var elemSearchText = document.getElementById('search_text');
   var elemMethod = document.getElementById('in_method');
   
   elemInSearch.value = elemSearchText.value;
   elemMethod.value = 'TpbFile';
   
   if (document.pressed == 'user')
   {
      elemMethod.value = "TpbUser";
   }
   
   return true;
}

var json_template = {};
var transform = {'tag':'li','html':'${Name} (${Category})'};
var data = [
    {'name':'Bob','age':40},
    {'name':'Frank','age':15},
    {'name':'Bill','age':65},
    {'name':'Robert','age':24}
];

var json_test = {"tag":"table","id":"SearchResult","children":[
    {"tag":"tbody","children":[
        {"tag":"tr","children":[
            {"tag":"td","children":[
                {"tag":"div","class":"detName","children":[
                    {"tag":"a","href":"${Link}","class":"detLink","title":"Details for ${Name}","html":"${Name}"}
                  ]},
                {"tag":"a","href":"${Magnet}","title":"Download this torrent using magnet","children":[
                    {"tag":"img","src":"icon-magnet.gif","alt":"Magnet link","html":""}
                  ]},
                {"tag":"a","href":"${TorrentLink}","title":"Download this torrent","children":[
                    {"tag":"img","src":"dl.gif","class":"dl","alt":"Download","html":""}
                  ]},
                {"tag":"img","src":"icon_comment.gif","alt":"This torrent has 1 comments.","title":"This torrent has 1 comments.","html":""},
                {"tag":"img","src":"icon_image.gif","alt":"This torrent has a cover image","title":"This torrent has a cover image","html":""},
                {"tag":"font","class":"detDesc","children":[
                    {"tag":"a","class":"detDesc","href":"${UserLink}","title":"Browse ${UserName}","html":"${UserName}"}
                  ]}
              ]},
            {"tag":"td","align":"right","html":"${Seed}"},
            {"tag":"td","align":"right","html":"${Leech}"}
          ]}
      ]}
  ]}

var xmlHttp = null;

function ProcessRequest() 
{
    if ( xmlHttp.readyState == 4 && xmlHttp.status == 200 ) 
    {
		var test_div = document.getElementById('test_div');
		
		var info = JSON.parse(xmlHttp.responseText);
		console.log(info);
		
        if ( info == null || info.Results == null || !$.isArray(info.Results) || info.Results.length == 0 ) 
        {
			test_div.innerHTML = "No results found";
        }
        else
        {
			test_div.innerHTML = json2html.transform(info.Results, json_test);
        }
    }
}

var HttpClient = function() {
    this.get = function(theUrl, aCallback) {
        xmlHttp = new XMLHttpRequest();
        xmlHttp.onreadystatechange = function() { 
            if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
                aCallback(xmlHttp.responseText);
        }

        xmlHttp.open( "GET", theUrl, true );            
        xmlHttp.send( null );
    }
}


/** Testet with:
 *  - IE 5.5, 7.0, 8.0, 9.0 (preview)
 *  - Firefox 3.6.3, 3.6.8
 *  - Safari 5.0
 *  - Chrome 5.0
 *  - Opera 10.10, 10.60
 */
var JavaScript = {
  load: function(src, callback) {
    var script = document.createElement('script'),
        loaded;
    script.setAttribute('src', src);
    if (callback) {
      script.onreadystatechange = script.onload = function() {
        if (!loaded) {
          callback();
        }
        loaded = true;
      };
    }
    document.getElementsByTagName('head')[0].appendChild(script);
  }
};


// Displays the request string with the colors and execute the request
function doRequest(form) {
   var elems = form.elements;
   //var iframe = document.getElementById('xmlOutputFrame');
   var querySpan = document.getElementById('query');
   var requestParams = [];
   var formattedRequestString = '';
   var value, name, requestString;
   var i;
   var test_div = document.getElementById('test_div');
   var test_json;
   //var itest = document.getElementById('itest');
   
   doHandleSubmit(form);
   
   //iframe.src = "about:blank";
   for (i = 0; i != elems.length; i++) {
      if (!(name = elems[i].name) || name == '_environment' || name == '_autofill') {
         continue;
      }

      if (elems[i].type == 'text' || elems[i].type == 'hidden' || elems[i].type == 'textarea') {
         value = elems[i].value;
      } else if (elems[i].type == 'select-one') {
         value = elems[i].options[elems[i].selectedIndex].value;
      }

      if (value) {
         if (name == '_action' || name == '_method' || name == '_target') {
            name = name.substring(1);
         }
         if (window.encodeURIComponent) {
            value = encodeURIComponent(value);
         } else {
            value = escape(value);
         }
         requestParams[requestParams.length] = name + '=' + value;
         if (formattedRequestString) {
            formattedRequestString += '&amp;';
         }

         if (name == 'method') {
            formattedRequestString += '<span class="functionparam">';
         } else {
            formattedRequestString += '<span class="param">';
         }

         formattedRequestString += '<span class="name">' + name + '</span>';
         formattedRequestString += '=<span class="value">' + value + '</span>';
         formattedRequestString += '</span>';
      }
   }

   requestString = form.action + '?' + requestParams.join('&');
   formattedRequestString = form.action + '?' + formattedRequestString;

   querySpan.innerHTML = formattedRequestString;
   
   test_div.innerHTML = "";
   console.log(requestString);
   
   clientReq = new HttpClient();
   clientReq.get(requestString, function(responseText) {
		var info = JSON.parse(responseText);
		console.log(info);
        if ( info == null || info.Results == null || !$.isArray(info.Results) || info.Results.length == 0 ) 
        {
		   test_div.innerHTML = "No results found";
        }
        else
        {
		   test_div.innerHTML = json2html.transform(info.Results, json_test);
        }
   });
   
   return false;
}

function getCookie(name) {
   var start = document.cookie.indexOf(name + "=");
   var len = start + name.length + 1;
   if ((!start) && (name != document.cookie.substring(0, name.length))) {
      return null;
   }
   if (start == -1) { return null; }
   var end = document.cookie.indexOf(";", len);
   if (end == -1) { end = document.cookie.length; }
   return unescape(document.cookie.substring(len, end));
}

function setCookie(name, value, expires, path, domain, secure) {
   var today = new Date();
   today.setTime(today.getTime());
   if (expires) {
      expires = expires * 1000 * 60 * 60 * 24;
   }
   var expires_date = new Date(today.getTime() + (expires));
   document.cookie = name + "=" + escape(value) +
      ((expires) ? ";expires=" + expires_date.toGMTString() : "") + //expires.toGMTString()
      ((path) ? ";path=" + path : "") +
      ((domain) ? ";domain=" + domain : "") +
      ((secure) ? ";secure" : "");
}


function deleteCookie(name, path, domain) {
   if (getCookie(name)) { document.cookie = name + "=" +
      ((path) ? ";path=" + path : "") +
      ((domain) ? ";domain=" + domain : "") +
      ";expires=Thu, 01-Jan-1970 00:00:01 GMT";
   }
}


function setEnvCookie(form) {
   var env;
   if (form._environment.options) {
      var selIndex = form._environment.selectedIndex;
      env = form._environment.options[selIndex].text;
      setCookie("xins.env", env, "", "", "", "");
   } else { //if (form._environment.type == 'text') {
      env = form._environment.value;
      setCookie("xins.env", env, "", "", "", "");
   }
}


function selectEnv() {
   var i;
   // make sure that only pages with form and environment set selected value from the env cookie
   if (document.forms[0] && document.forms[0]._environment && document.forms[0]._environment.options) {
      var options = document.forms[0]._environment.options;
      env = getCookie("xins.env");
      for (i = 0; i != options.length; i++) {
         var option = options[i];
         if (env == options[i].text) {
            options.selectedIndex = i;
         }
      }
   } else if (document.forms[0] && document.forms[0]._environment && document.forms[0]._environment.type == 'text') {
      env = getCookie("xins.env");
      if (env != null && env != 'null') {
         document.forms[0]._environment.value = env;
      }
   }
}
