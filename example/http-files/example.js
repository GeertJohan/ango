var chat = angular.module('chatExample', ['ango-chatService']);

chat.config(function(chatServiceProvider) {
	chatServiceProvider.setDebug(true);
	var fn = function(err) {
		console.error("catched err! " + err);
	};
	chatServiceProvider.listenOnWsError(fn)
});

chat.controller('mainCtrl', function($scope, chatService) {
	$scope.foo = "bar";
	console.log(chatService.getServiceName());
	console.log(chatService.getProtocolVersion());
});