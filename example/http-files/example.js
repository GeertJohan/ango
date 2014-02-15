var app = angular.module('app', ['chatService']);

app.config(function(chatServiceProvider) {
	chatServiceProvider.setUri("");
	chatServiceProvider.handle('newMessage', function(message) {
		console.log('Have a new message');
	});
});

app.controller('chatCtrl', function($scope, chatService) {
	$scope.foo = "bar";
})