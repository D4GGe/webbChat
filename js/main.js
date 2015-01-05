var webChat = angular.module('webChat', ['ngRoute']);

webChat.config(['$routeProvider',
  function ($routeProvider) {
      $routeProvider.
        when('/', {
            templateUrl: 'views/chatView.html',
            controller: 'ChatCtrl'
        }).
        otherwise({
            redirectTo: '/'
        });
  }]);



webChat.controller('ChatCtrl', function ($scope, $http) {
    $scope.posts = [];

    $http.post('/', { msg: 'hello word!' }).
  success(function (data, status, headers, config) {
      // this callback will be called asynchronously
      // when the response is available
  }).
  error(function (data, status, headers, config) {
      // called asynchronously if an error occurs
      // or server returns response with an error status.
  });


    $scope.updatePosts = function () {
        var i = 0;
        if ($scope.posts.length != 0)
            i = $scope.posts[$scope.posts.length - 1].Id;

        $http.get('http://192.168.1.200:8080/chat/hello?id='+i).
      success(function (data, status, headers, config) {
          $scope.posts = $scope.posts.concat(data);

      }).
      error(function (data, status, headers, config) {
          // called asynchronously if an error occurs
          // or server returns response with an error status.
      })
    };


    $scope.sendPost = function () {
        $http.get('http://192.168.1.200:8080/chat/hello?msg=' + $scope.sendText).
         success(function (data, status, headers, config) {
             $scope.updatePosts();
             $scope.sendText="";
         }).
         error(function (data, status, headers, config) {
             // called asynchronously if an error occurs
             // or server returns response with an error status.
         })
    }

    setInterval($scope.updatePosts, 100);
    
});