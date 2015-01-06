var webChat = angular.module('webChat', ['ngRoute']);

webChat.config(['$routeProvider',
  function ($routeProvider) {
      $routeProvider.
        when('/room/:roomName', {
            templateUrl: 'views/chatView.html',
            controller: 'ChatCtrl'
        }).
        otherwise({
            redirectTo: '/room/hello'
        });
  }]);



webChat.controller('ChatCtrl', function ($scope, $http, $routeParams) {
    $scope.posts = [];
    $scope.roomName = $routeParams.roomName;

    $scope.updatePosts = function () {
        var i = 0;
        if ($scope.posts.length != 0)
            i = $scope.posts[$scope.posts.length - 1].Id;

        $http.get('http://192.168.1.200:8080/chat/'+ $scope.roomName +'/'+i).
      success(function (data, status, headers, config) {
          $scope.posts = $scope.posts.concat(data);

      }).
      error(function (data, status, headers, config) {
          // called asynchronously if an error occurs
          // or server returns response with an error status.
      })
    };


    $scope.sendPost = function () {
    	alert($scope.sendText);
        $http.post('http://192.168.1.200:8080/chat/' + $scope.roomName, {msg:$scope.sendText}).
         success(function (data, status, headers, config) {
             $scope.updatePosts();
             $scope.sendText="";
         }).
         error(function (data, status, headers, config) {
             // called asynchronously if an error occurs
             // or server returns response with an error status.
         })
    }

    setInterval($scope.updatePosts, 2000);
    
});