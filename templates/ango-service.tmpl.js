angular.module('{{.ServiceName}}', [])
.provider('{{.ServiceName}}', function($q) {
	var provider = {};

	provider.setUrl = function(url) {
		provider.url = url;
	};

	provider.$get = function() {
		//++
	};
})
