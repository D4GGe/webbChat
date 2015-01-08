var webChat = angular.module('webChat', ['ngRoute']);

webChat.config(['$routeProvider',
  function ($routeProvider) {
      $routeProvider.
        when('/room/:roomName', {
            templateUrl: 'views/chatViewNew.html',
            controller: 'ChatCtrl'
        }).
        otherwise({
            redirectTo: '/room/hello'
        });
  }]);

webChat.controller('ChatCtrl', function ($scope, $http, $routeParams, $rootScope) {
    if ($rootScope.posts == null)
        $rootScope.posts = [];
    $scope.roomName = $routeParams.roomName;
    if ($rootScope.posts[$scope.roomName] == null)
        $rootScope.posts[$scope.roomName] = [];


    $scope.posts = $rootScope.posts[$scope.roomName];


    

    $scope.updatePosts = function () {
        var i = 0;
        if ($rootScope.posts[$scope.roomName].length != 0)
            i = $rootScope.posts[$scope.roomName][$rootScope.posts[$scope.roomName].length - 1].Id;

        $http.get('http://192.168.1.200:8080/chat/'+ $scope.roomName +'/'+i).
      success(function (data, status, headers, config) {
          var elem = document.getElementById('scrollBody');
         
          $rootScope.posts[$scope.roomName] = $rootScope.posts[$scope.roomName].concat(data);
          $scope.posts = $rootScope.posts[$scope.roomName];

          if (elem.scrollHeight - elem.scrollTop < 1000 && data.length != 0)
              setTimeout(function () {
                  var elem = document.getElementById('scrollBody');
                  elem.scrollTop = elem.scrollHeight;
              }, 30)

      }).
      error(function (data, status, headers, config) {
          // called asynchronously if an error occurs
          // or server returns response with an error status.
      })
    };

    $scope.keyPress = function (e) {
        if (e.which == 13)
            $scope.sendPost();
    }

    $scope.sendPost = function () {
    	var req = {
		 method: 'POST',
		 url: 'http://192.168.1.200:8080/chat/' + $scope.roomName,
		 headers: {
		   'Content-Type': 'application/x-www-form-urlencoded'
		 },
		 data: "msg="+$scope.sendText,
		};



        $http(req).
         success(function (data, status, headers, config) {
             //$scope.updatePosts(); Quick fix
             $scope.sendText = "";
         }).
         error(function (data, status, headers, config) {
             // called asynchronously if an error occurs
             // or server returns response with an error status.
         })
    }

    setInterval($scope.updatePosts, 500);
    
});